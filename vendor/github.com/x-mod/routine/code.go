package routine

import (
	"syscall"
)

//Reference: https://www.tldp.org/LDP/abs/html/exitcodes.html

//Code for process exit
type Code int32

//go:generate stringer -type Code code.go
const (
	// OK is returned on success.
	OK              Code = 0
	GeneralErr      Code = 1
	MisUseErr       Code = 2
	NotExecutable   Code = 126
	IllegalCommand  Code = 127
	InvalidArgments Code = 128
	// SIGNAL exits = 128 + SIGNAL
)

//SignalCode signal code
func SignalCode(sig syscall.Signal) Code {
	return Code(128 + int(sig))
}

//Value of Code
func (c Code) Value() int32 {
	return int32(c)
}
