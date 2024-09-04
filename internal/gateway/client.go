package gateway

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/stringutil"
)

type Client struct {
	ctx        *ConnContext
	conn       socketio.Conn
	PlatformID int
	UserID     string
	token      string
}

func (c *Client) ResetClient(ctx *ConnContext, conn socketio.Conn) {
	c.ctx = ctx
	c.conn = conn
	c.PlatformID = stringutil.StringToInt(ctx.GetPlatformID())
	c.UserID = ctx.GetUserID()
	c.token = ctx.GetToken()
}

func (c *Client) KickOnlineMessage() error {
	resp := Resp{
		ReqIdentifier: WSKickOnlineMsg,
	}
	log.ZInfo(c.ctx, "KickOnlineMessage", resp.String())
	c.conn.Emit(SocketResponseEvent, resp.String())
	return c.conn.Close()
}
