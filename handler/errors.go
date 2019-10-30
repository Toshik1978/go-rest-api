package handler

import (
	"errors"

	"github.com/Toshik1978/go-rest-api/service/errutil"
)

// ErrorKind define kind of error, occurred
type ErrorKind int8

const (
	ServerError ErrorKind = iota + 1
	ClientError
)

// Error define custom handler error
type Error struct {
	error
	Kind ErrorKind
}

// NewError creates new custom handler error
func NewError(message string, kind ErrorKind) error {
	return &Error{
		error: errors.New(message),
		Kind:  kind,
	}
}

// WrapError creates new wrapped handler error
func WrapError(err error, message string, kind ErrorKind) error {
	return &Error{
		error: errutil.Wrap(err, message),
		Kind:  kind,
	}
}
