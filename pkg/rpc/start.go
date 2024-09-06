package rpc

import (
	"context"
	"fmt"
	"github.com/liserc/open-socket/pkg/config"
	"github.com/liserc/open-socket/pkg/prommetrics"
	"github.com/openimsdk/tools/mw"
	"github.com/openimsdk/tools/utils/datautil"
	"google.golang.org/grpc/status"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/system/program"
	"github.com/openimsdk/tools/utils/network"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Start rpc server.
func Start[T any](ctx context.Context, discovery *config.Discovery, prometheusConfig *config.Prometheus, listenIP,
	registerIP string, rpcPorts []int, index int, rpcRegisterName string, share *config.Share, config T, rpcFn func(ctx context.Context,
	config T, client discovery.SvcDiscoveryRegistry, server *grpc.Server) error, options ...grpc.ServerOption) error {

	rpcPort, err := datautil.GetElemByIndex(rpcPorts, index)
	if err != nil {
		return err
	}

	log.CInfo(ctx, "RPC server is initializing", "rpcRegisterName", rpcRegisterName, "rpcPort", rpcPort,
		"prometheusPorts", prometheusConfig.Ports)
	rpcTcpAddr := net.JoinHostPort(network.GetListenIP(listenIP), strconv.Itoa(rpcPort))
	listener, err := net.Listen(
		"tcp",
		rpcTcpAddr,
	)
	if err != nil {
		return errs.WrapMsg(err, "listen err", "rpcTcpAddr", rpcTcpAddr)
	}

	defer listener.Close()
	client, err := NewDiscoveryRegister(discovery, share)
	if err != nil {
		return err
	}

	defer client.Close()
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	registerIP, err = network.GetRpcRegisterIP(registerIP)
	if err != nil {
		return err
	}

	if prometheusConfig.Enable {
		options = append(
			options, mw.GrpcServer(),
			prommetricsUnaryInterceptor(rpcRegisterName),
			prommetricsStreamInterceptor(rpcRegisterName),
		)
	} else {
		options = append(options, mw.GrpcServer())
	}

	srv := grpc.NewServer(options...)
	once := sync.Once{}
	defer func() {
		once.Do(srv.GracefulStop)
	}()

	err = rpcFn(ctx, config, client, srv)
	if err != nil {
		return err
	}

	err = client.Register(
		rpcRegisterName,
		registerIP,
		rpcPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	var (
		netDone    = make(chan struct{}, 2)
		netErr     error
		httpServer *http.Server
	)
	if prometheusConfig.Enable {
		go func() {
			prometheusPort, err := datautil.GetElemByIndex(prometheusConfig.Ports, index)
			if err != nil {
				netErr = err
				netDone <- struct{}{}
				return
			}
			cs := prommetrics.GetGrpcCusMetrics(rpcRegisterName, share)
			if err := prommetrics.RpcInit(cs, prometheusPort); err != nil && err != http.ErrServerClosed {
				netErr = errs.WrapMsg(err, fmt.Sprintf("rpc %s prometheus start err: %d", rpcRegisterName, prometheusPort))
				netDone <- struct{}{}
			}
		}()
	}

	go func() {
		err = srv.Serve(listener)
		if err != nil {
			netErr = errs.WrapMsg(err, "rpc start err: ", rpcTcpAddr)
			netDone <- struct{}{}
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	select {
	case <-sigs:
		program.SIGTERMExit()
		back, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err = gracefulStopWithCtx(back, srv.GracefulStop); err != nil {
			return err
		}
		back, cancel = context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		err = httpServer.Shutdown(back)
		if err != nil {
			return errs.WrapMsg(err, "shutdown err")
		}
		return nil
	case <-netDone:
		return netErr
	}
}

func gracefulStopWithCtx(ctx context.Context, f func()) error {
	done := make(chan struct{}, 1)
	go func() {
		f()
		close(done)
	}()
	select {
	case <-ctx.Done():
		return errs.New("timeout, ctx graceful stop")
	case <-done:
		return nil
	}
}

func prommetricsUnaryInterceptor(rpcRegisterName string) grpc.ServerOption {
	getCode := func(err error) int {
		if err == nil {
			return 0
		}
		rpcErr, ok := err.(interface{ GRPCStatus() *status.Status })
		if !ok {
			return -1
		}
		return int(rpcErr.GRPCStatus().Code())
	}
	return grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		prommetrics.RPCCall(rpcRegisterName, info.FullMethod, getCode(err))
		return resp, err
	})
}

func prommetricsStreamInterceptor(rpcRegisterName string) grpc.ServerOption {
	return grpc.ChainStreamInterceptor()
}
