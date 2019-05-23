package errors

import (
	"errors"
	"fmt"

	stackerr "github.com/pkg/errors"
)

type causer interface {
	Cause() error
}

type coder interface {
	Value() uint32
}

//New errorstring error
func New(err string) error {
	return errors.New(err)
}

//Errorf standard func
func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

//WithStack error
func WithStack(err error) error {
	return stackerr.WithStack(err)
}
