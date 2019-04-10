package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//Code interface
type Code interface {
	Value() int32
	String() string
}

type codeErr struct {
	code Code
	err  error
}

//WithCode wrap err with code
func WithCode(err error, code Code) error {
	if err != nil {
		return &codeErr{
			err:  err,
			code: code,
		}
	}
	return nil
}

//CodeError a new error from code
func CodeError(code Code) error {
	if code != nil {
		return WithCode(New(code.String()), code)
	}
	return nil
}

func (ce *codeErr) Error() string {
	return fmt.Sprintf("%d: %s", ce.code.Value(), ce.err.Error())
}

func (ce *codeErr) Cause() error {
	return ce.err
}

func (ce *codeErr) Value() int32 {
	if ce.code != nil {
		return ce.code.Value()
	}
	return 0
}

//GRPCStatus make codeErr support grpc status
func (ce *codeErr) GRPCStatus() *status.Status {
	return status.New(codes.Code(ce.code.Value()), ce.err.Error())
}
