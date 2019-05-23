package errors

import (
	"fmt"
)

//internal Code implemention
type errorCode struct {
	value   int32
	message string
}

func (code *errorCode) Value() int32 {
	return code.value
}

func (code *errorCode) String() string {
	if code.message != "" {
		return code.message
	}
	return fmt.Sprintf("Error(%d)", code.value)
}

//Error implement error
func (code *errorCode) Error() string {
	return code.String()
}

//ValueErr pure value error
func ValueErr(v int32) error {
	if v != 0 {
		return &errorCode{value: v}
	}
	return nil
}
