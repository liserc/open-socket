package gateway

import "github.com/openimsdk/tools/utils/jsonutil"

type Req struct {
	ReqIdentifier int32  `json:"cmd"`
	Token         string `json:"token"`
	SendID        string `json:"sender"`
	OperationID   string `json:"operation"`
	Data          []byte `json:"data"`
}

func (r *Req) String() string {
	var req Req
	req.ReqIdentifier = r.ReqIdentifier
	req.Token = r.Token
	req.SendID = r.SendID
	req.OperationID = r.OperationID
	req.Data = r.Data
	return jsonutil.StructToJsonString(req)
}

type Resp struct {
	ReqIdentifier int32  `json:"cmd"`
	OperationID   string `json:"operation"`
	ErrCode       int    `json:"code"`
	ErrMsg        string `json:"message"`
	Data          string `json:"data"`
}

func (r *Resp) String() string {
	var resp Resp
	resp.ReqIdentifier = r.ReqIdentifier
	resp.OperationID = r.OperationID
	resp.ErrCode = r.ErrCode
	resp.ErrMsg = r.ErrMsg
	resp.Data = r.Data
	return jsonutil.StructToJsonString(resp)
}
