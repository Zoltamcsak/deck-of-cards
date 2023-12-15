package error

import (
	"errors"
)

type (
	// Error custom error type
	Error struct {
		kind    Kind
		message string
		cause   error
	}

	Kind int
)

func New(kind Kind, message string) *Error {
	return &Error{
		kind:    kind,
		message: message,
		cause:   errors.New(message),
	}
}

func Wrap(kind Kind, message string, err error) *Error {
	return &Error{
		kind:    kind,
		message: message,
		cause:   err,
	}
}

func (e *Error) Error() string {
	if e.cause != nil {
		return e.cause.Error()
	}
	return e.message
}

func (e *Error) Kind() int {
	return int(e.kind)
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Cause() error {
	return e.cause
}
