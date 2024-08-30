package errs

var (
	ErrArgs             = NewCodeError(ArgsError, "ArgsError")
	ErrNoPermission     = NewCodeError(NoPermissionError, "NoPermissionError")
	ErrInternalServer   = NewCodeError(ServerInternalError, "ServerInternalError")
	ErrRecordNotFound   = NewCodeError(RecordNotFoundError, "RecordNotFoundError")
	ErrDuplicateKey     = NewCodeError(DuplicateKeyError, "DuplicateKeyError")
	ErrTokenExpired     = NewCodeError(TokenExpiredError, "TokenExpiredError")
	ErrTokenInvalid     = NewCodeError(TokenInvalidError, "TokenInvalidError")
	ErrTokenMalformed   = NewCodeError(TokenMalformedError, "TokenMalformedError")
	ErrTokenNotValidYet = NewCodeError(TokenNotValidYetError, "TokenNotValidYetError")
	ErrTokenUnknown     = NewCodeError(TokenUnknownError, "TokenUnknownError")
	ErrTokenKicked      = NewCodeError(TokenKickedError, "TokenKickedError")
	ErrTokenNotExist    = NewCodeError(TokenNotExistError, "TokenNotExistError")
)
