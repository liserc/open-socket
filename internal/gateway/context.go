package gateway

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/liserc/open-socket/pkg/servererrs"
	"github.com/openimsdk/protocol/constant"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/openimsdk/tools/utils/stringutil"
)

type ConnContext struct {
	URL          url.URL
	Path         string
	RemoteHeader http.Header
	RemoteAddr   string
	ConnID       string
	ConnErr      error
}

func (c *ConnContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *ConnContext) Done() <-chan struct{} {
	return nil
}

func (c *ConnContext) Err() error {
	return nil
}

func (c *ConnContext) Value(key any) any {
	switch key {
	case constant.OpUserID:
		return c.GetUserID()
	case constant.OperationID:
		return c.GetOperationID()
	case constant.ConnID:
		return c.GetConnID()
	case constant.OpUserPlatform:
		return constant.PlatformIDToName(stringutil.StringToInt(c.GetPlatformID()))
	case constant.RemoteAddr:
		return c.RemoteAddr
	default:
		return ""
	}
}

func newConnContext(conn socketio.Conn) *ConnContext {
	return &ConnContext{
		URL:        conn.URL(),
		Path:       conn.URL().Path,
		RemoteAddr: conn.RemoteAddr().String(),
		ConnID:     conn.ID(),
	}
}

func (c *ConnContext) GetRemoteAddr() string {
	return c.RemoteAddr
}

func (c *ConnContext) Query(key string) (string, bool) {
	var value string
	if value = c.URL.Query().Get(key); value == "" {
		return value, false
	}
	return value, true
}

func (c *ConnContext) GetHeader(key string) (string, bool) {
	var value string
	if value = c.RemoteHeader.Get(key); value == "" {
		return value, false
	}
	return value, true
}

func (c *ConnContext) GetConnID() string {
	return c.ConnID
}

func (c *ConnContext) GetUserID() string {
	return c.URL.Query().Get(WsUserID)
}

func (c *ConnContext) GetPlatformID() string {
	return c.URL.Query().Get(PlatformID)
}

func (c *ConnContext) GetOperationID() string {
	return c.URL.Query().Get(OperationID)
}

func (c *ConnContext) SetOperationID(operationID string) {
	values := c.URL.Query()
	values.Set(OperationID, operationID)
	c.URL.RawQuery = values.Encode()
}

func (c *ConnContext) GetToken() string {
	return c.URL.Query().Get(Token)
}

func (c *ConnContext) GetCompression() bool {
	compression, exists := c.Query(Compression)
	if exists && compression == GzipCompressionProtocol {
		return true
	} else {
		compression, exists := c.GetHeader(Compression)
		if exists && compression == GzipCompressionProtocol {
			return true
		}
	}
	return false
}

func (c *ConnContext) ShouldSendResp() bool {
	errResp, exists := c.Query(SendResponse)
	if exists {
		b, err := strconv.ParseBool(errResp)
		if err != nil {
			return false
		} else {
			return b
		}
	}
	return false
}

func (c *ConnContext) SetToken(token string) {
	c.URL.RawQuery = Token + "=" + token
}

func (c *ConnContext) GetBackground() bool {
	b, err := strconv.ParseBool(c.URL.Query().Get(BackgroundStatus))
	if err != nil {
		return false
	}
	return b
}

func (c *ConnContext) ParseEssentialArgs() error {
	_, exists := c.Query(Token)
	if !exists {
		return servererrs.ErrConnArgsErr.WrapMsg("token is empty")
	}
	_, exists = c.Query(WsUserID)
	if !exists {
		return servererrs.ErrConnArgsErr.WrapMsg("sendID is empty")
	}
	platformIDStr, exists := c.Query(PlatformID)
	if !exists {
		return servererrs.ErrConnArgsErr.WrapMsg("platformID is empty")
	}
	_, err := strconv.Atoi(platformIDStr)
	if err != nil {
		return servererrs.ErrConnArgsErr.WrapMsg("platformID is not int")
	}
	return nil
}
