package servererrs

var (
	// Token error codes.
	TokenExpiredError     = 1501
	TokenInvalidError     = 1502
	TokenMalformedError   = 1503
	TokenNotValidYetError = 1504
	TokenUnknownError     = 1505
	TokenKickedError      = 1506
	TokenNotExistError    = 1507

	// Long connection gateway error codes.
	ConnOverMaxNumLimit = 1601
	ConnArgsErr         = 1602
	PushMsgErr          = 1603
)
