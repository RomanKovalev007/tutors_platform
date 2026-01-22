package models

type ErrorResponseErrorCode string

const (
	USERNOTFOUND       ErrorResponseErrorCode = "USER_NOT_FOUND"
	USEREXISTS         ErrorResponseErrorCode = "USER_EXISTS"
	INVALIDTOKEN       ErrorResponseErrorCode = "INVALID_TOKEN"
	TOKENEXPIRED       ErrorResponseErrorCode = "TOKEN_EXPIRED"
	INTERNALERROR      ErrorResponseErrorCode = "INTERNAL_ERROR"
	INVALIDINPUT       ErrorResponseErrorCode = "INVALID_INPUT"
	INVALIDCREDENTIALS ErrorResponseErrorCode = "INVALID_CREDENTIALS"
	STATUS_OK          ErrorResponseErrorCode = "STATUS_OK"
)

type Error struct {
	Code    ErrorResponseErrorCode
	Message error
}
