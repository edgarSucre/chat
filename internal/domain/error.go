package domain

import "fmt"

type (
	Err struct {
		orig error
		msg  string
		code ErrorCode
	}

	ErrorCode int
)

const (
	// Validation
	ErrorCodeInvalidParams ErrorCode = iota

	// Repository
	ErrorCodeUserNotFound
	ErrorCodeUserConflict
	ErrorCodeInternalRepository
)

func WrapErrorf(orig error, code ErrorCode, format string, a ...interface{}) *Err {
	return &Err{
		code: code,
		orig: orig,
		msg:  fmt.Sprintf(format, a...),
	}
}

func (e *Err) Error() string {
	if e == nil {
		return ""
	}

	if e.orig != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.orig)
	}

	return e.msg
}

func (e *Err) Unwrap() error {
	return e.orig
}

func (e *Err) Msg() string {
	return e.msg
}

func (e *Err) Code() ErrorCode {
	return e.code
}
