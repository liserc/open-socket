package gateway

import (
	"context"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/liserc/open-socket/pkg/prommetrics"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/stringutil"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type LongConnServer interface {
	Run(done chan error) error
	SetDiscoveryRegistry(client discovery.SvcDiscoveryRegistry)
	GetUserAllCons(userID string) ([]*Client, bool)
	GetUserPlatformCons(userID string, platform int) ([]*Client, bool, bool)
	KickUserConn(client *Client) error
}

type SocketServer struct {
	config            *Config
	port              int
	maxConnNum        int64
	onlineUserNum     atomic.Int64
	onlineUserConnNum atomic.Int64
	registerChan      chan *Client
	unregisterChan    chan *Client
	clients           UserMap
	clientPool        sync.Pool
	socket            *socketio.Server
	registry          discovery.SvcDiscoveryRegistry
}

func NewSocketServer(config *Config, opts ...Option) *SocketServer {
	var opt configs
	for _, o := range opts {
		o(&opt)
	}
	return &SocketServer{
		config:     config,
		port:       opt.port,
		maxConnNum: opt.maxConnNum,
		clients:    newUserMap(),
		clientPool: sync.Pool{
			New: func() any {
				return new(Client)
			},
		},
	}
}

func (ws *SocketServer) Run(done chan error) error {
	var (
		err  error
		nete chan error
	)

	socket := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&websocket.Transport{
				CheckOrigin: handleAllowOrigin,
			},
			&polling.Transport{
				CheckOrigin: handleAllowOrigin,
			},
		},
	})
	socket.OnConnect("/", func(conn socketio.Conn) error {
		log.CInfo(context.Background(), "connect", "ID", conn.ID(), "RawQuery", conn.URL().RawQuery)
		ctx := newConnContext(conn)
		err = ctx.ParseEssentialArgs()
		if err != nil {
			log.CInfo(ctx, "connect args error", err)
			ctx.ConnErr = err
			conn.SetContext(ctx)
			return err
		}
		conn.SetContext(ctx)

		// Retrieve a client object from the client pool, reset its state, and associate it with the current WebSocket long connection
		client := ws.clientPool.Get().(*Client)
		client.ResetClient(ctx, conn)
		ws.registerClient(client)

		return nil
	})
	socket.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		ctx := conn.Context().(*ConnContext)
		log.CInfo(ctx, "disconnect", "reason", reason)
		if ctx.ConnErr != nil {
			log.ZError(ctx, "connect error close", ctx.ConnErr)
			_ = conn.Close()
			return
		}

		clients, _, ok := ws.GetUserPlatformCons(ctx.GetUserID(), stringutil.StringToInt(ctx.GetPlatformID()))
		if !ok {
			log.ZDebug(ctx, "conn not exist", "userID", ctx.GetUserID(), "platformID", ctx.GetPlatformID())
			return
		}
		for _, client := range clients {
			if client.conn.ID() != conn.ID() {
				continue
			}
			ws.unregisterClient(client)
		}
	})
	socket.OnError("/", func(conn socketio.Conn, e error) {
		ctx := conn.Context().(*ConnContext)
		log.CInfo(ctx, "error", "ConnErr", ctx.ConnErr)
	})
	socket.OnEvent("/", SocketRequestEvent, func(conn socketio.Conn, message string) {
		ctx := conn.Context().(*ConnContext)
		log.CInfo(ctx, SocketRequestEvent, "message", message)
	})
	ws.socket = socket

	nete = make(chan error)
	server := &http.Server{
		Addr:    ":" + stringutil.IntToString(ws.port),
		Handler: nil,
	}
	go func() {
		go func() {
			if err = socket.Serve(); err != nil {
				nete <- errs.WrapMsg(err, "socketio listen error")
			}
		}()
		http.Handle("/socket.io/", socket)
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			nete <- errs.WrapMsg(err, "server listen error", server.Addr)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	select {
	case err = <-done:
		closed := socket.Close()
		if closed != nil {
			return errs.WrapMsg(closed, "socket shutdown error")
		}
		down := server.Shutdown(ctx)
		if down != nil {
			return errs.WrapMsg(down, "server shutdown error")
		}

	case err = <-nete:
		closed := socket.Close()
		if closed != nil {
			return errs.WrapMsg(closed, "socket shutdown error")
		}
		down := server.Shutdown(ctx)
		if down != nil {
			return errs.WrapMsg(down, "server shutdown error")
		}
	}

	return err
}

func (ws *SocketServer) SetDiscoveryRegistry(registry discovery.SvcDiscoveryRegistry) {
	ws.registry = registry
}

func (ws *SocketServer) GetUserAllCons(userID string) ([]*Client, bool) {
	return ws.clients.GetAll(userID)
}

func (ws *SocketServer) GetUserPlatformCons(userID string, platform int) ([]*Client, bool, bool) {
	return ws.clients.Get(userID, platform)
}

func (ws *SocketServer) KickUserConn(client *Client) error {
	ws.clients.DeleteClients(client.UserID, []*Client{client})
	return client.KickOnlineMessage()
}

func (ws *SocketServer) registerClient(client *Client) {
	var (
		userOK     bool
		clientOK   bool
		oldClients []*Client
	)
	oldClients, userOK, clientOK = ws.clients.Get(client.UserID, client.PlatformID)
	if !userOK {
		ws.clients.Set(client.UserID, client)
		log.ZDebug(client.ctx, "user not exist", "userID", client.UserID, "platformID", client.PlatformID)
		prommetrics.OnlineUserGauge.Add(1)
		ws.onlineUserNum.Add(1)
		ws.onlineUserConnNum.Add(1)
	} else {
		log.ZDebug(client.ctx, "user exist", "userID", client.UserID, "platformID", client.PlatformID)
		if clientOK {
			ws.clients.Set(client.UserID, client)
			// There is already a connection to the platform
			log.ZInfo(client.ctx, "repeat login", "userID", client.UserID, "platformID", client.PlatformID, "old remote addr", getRemoteAdders(oldClients))
			ws.onlineUserConnNum.Add(1)
		} else {
			ws.clients.Set(client.UserID, client)
			ws.onlineUserConnNum.Add(1)
		}
	}

	wg := sync.WaitGroup{}
	wg.Wait()
	log.ZInfo(client.ctx, "user online", "online user nums", ws.onlineUserNum.Load(), "online user conn nums", ws.onlineUserConnNum.Load())
}

func (ws *SocketServer) unregisterClient(client *Client) {
	defer ws.clientPool.Put(client)
	isDeleteUser := ws.clients.DeleteClients(client.UserID, []*Client{client})
	if isDeleteUser {
		ws.onlineUserNum.Add(-1)
		prommetrics.OnlineUserGauge.Dec()
	}
	ws.onlineUserConnNum.Add(-1)
	log.ZInfo(client.ctx, "user offline", "close conn", client.conn.ID(),
		"online user num", ws.onlineUserNum.Load(), "online user conn num", ws.onlineUserConnNum.Load(),
	)
}

func handleAllowOrigin(r *http.Request) bool {
	log.CInfo(context.Background(), "Origin Request", "RequestURI", r.RequestURI, "Header", r.Header)
	return true
}

func getRemoteAdders(client []*Client) string {
	var ret string
	for i, c := range client {
		if i == 0 {
			ret = c.ctx.GetRemoteAddr()
		} else {
			ret += "@" + c.ctx.GetRemoteAddr()
		}
	}
	return ret
}
