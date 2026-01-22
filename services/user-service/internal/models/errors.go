package models

type ErrorResponseErrorCode string

const (
	USERNOTFOUND    ErrorResponseErrorCode = "USER_NOT_FOUND"
	TUTORNOTFOUND   ErrorResponseErrorCode = "TUTOR_NOT_FOUND"
	STUDENTNOTFOUND ErrorResponseErrorCode = "STUDENT_NOT_FOUND"
	USEREXISTS      ErrorResponseErrorCode = "USER_EXISTS"
	TUTOREXISTS     ErrorResponseErrorCode = "TUTOR_EXISTS"
	STUDENTEXISTS   ErrorResponseErrorCode = "STUDENT_EXISTS"
	INTERNALERROR   ErrorResponseErrorCode = "INTERNAL_ERROR"
	INVALIDINPUT    ErrorResponseErrorCode = "INVALID_INPUT"
	STATUS_OK       ErrorResponseErrorCode = "STATUS_OK"
)

type Error struct {
	Code    ErrorResponseErrorCode
	Message error
}
