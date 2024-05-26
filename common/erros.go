package common

import (
	"errors"
	"net/http"
)

var (
	RecordNotFound               = errors.New("record not found")
	ErrNameCannotBeEmpty         = errors.New("name cannot be empty")
	ErrUserNameOrPasswordInvalid = AppError{
		StatusCode: http.StatusUnauthorized,
		RootErr:    errors.New("username or password invalid"),
		Message:    "Username or password invalid",
		Key:        "ErrUserNameOrPasswordInvalid",
	}
)
