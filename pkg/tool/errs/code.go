package errs

const (
	// General error codes.
	ServerInternalError = 500  // Server internal error
	ArgsError           = 1001 // Input parameter error
	NoPermissionError   = 1002 // Insufficient permission
	DuplicateKeyError   = 1003
	RecordNotFoundError = 1004 // Record does not exist

	TokenExpiredError     = 1501
	TokenInvalidError     = 1502
	TokenMalformedError   = 1503
	TokenNotValidYetError = 1504
	TokenUnknownError     = 1505
	TokenKickedError      = 1506
	TokenNotExistError    = 1507
)
