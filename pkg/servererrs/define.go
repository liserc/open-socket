package servererrs

import "github.com/openimsdk/tools/errs"

var (
	ErrConnArgsErr = errs.NewCodeError(ConnArgsErr, "args err, need token, sendID, platformID")
)
