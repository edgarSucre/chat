package domain

import "errors"

type (
	UserParam struct {
		UserName string
		Password string
	}

	UserResponse struct {
		ID       int64
		UserName string
		Password string
	}
)

var (
	ErrUserDoesNotExists   = errors.New("username entered does not exist")
	ErrWrongPassword       = errors.New("password is incorrect")
	ErrUserConflict        = errors.New("user already exists")
	ErrInternalServerError = errors.New("internal server error")
	ErrBadParamInput       = errors.New("given param are not valid")
)
