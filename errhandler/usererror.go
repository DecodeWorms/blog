package errhandler

import "net/http"

type UserError struct {
	Code int
	Msg  interface{}
}

func NewUserError(c int, m interface{}) *UserError {
	return &UserError{
		Code: c,
		Msg:  m,
	}
}

var (
	EmailAlreadyExist       = NewUserError(http.StatusConflict, "user with email already exists")
	PhoneNumberAlreadyExist = NewUserError(http.StatusConflict, "user with phone number already exists")
	UserRecordNotFound      = NewUserError(http.StatusNotFound, "user record not found")
	GenderNotMatched        = NewUserError(http.StatusBadRequest, "gender provided is not recognized")
	IncorrectPassword       = NewUserError(http.StatusInternalServerError, "the password is incorrect")
	ValidationErr           = NewUserError(http.StatusInternalServerError, "error in validation")
	UnableToCeateUser       = NewUserError(http.StatusInternalServerError, "unable to create a user")
	UnableToGenerateToken   = NewUserError(http.StatusInternalServerError, "unable to generate a token")
)
