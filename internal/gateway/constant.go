package gateway

import "time"

const (
	WsUserID                = "sendID"
	CommonUserID            = "userID"
	PlatformID              = "platformID"
	ConnID                  = "connID"
	Token                   = "token"
	OperationID             = "operationID"
	Compression             = "compression"
	GzipCompressionProtocol = "gzip"
	BackgroundStatus        = "background"
	SendResponse            = "sendResponse"
)

const (
	// Websocket Event.
	SocketRequestEvent  = "request"
	SocketResponseEvent = "response"
)

const (
	// Websocket Protocol.
	WSGetNewestSeq     = 1001
	WSPullMsgBySeqList = 1002
	WSSendMsg          = 1003
	WSSendSignalMsg    = 1004
	WSPushMsg          = 2001
	WSKickOnlineMsg    = 2002
	WsLogoutMsg        = 2003
	WSDataError        = 3001
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 51200
)
