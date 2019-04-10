package routine

import (
	"context"
	"syscall"
)

// InterruptHandler definition
type InterruptHandler func(ctx context.Context) (exit bool)

// Interruptor definition
type Interruptor interface {
	Signal() syscall.Signal
	Interrupt() InterruptHandler
}

//DefaultCancelInterruptors include INT/TERM/KILL signals
var DefaultCancelInterruptors []Interruptor

// CancelInterruptor definition
type CancelInterruptor struct {
	sig syscall.Signal
}

// NewCancelInterruptor if fn is nil will cancel context
func NewCancelInterruptor(sig syscall.Signal) *CancelInterruptor {
	return &CancelInterruptor{
		sig: sig,
	}
}

// Signal inplement the interface
func (c *CancelInterruptor) Signal() syscall.Signal {
	return c.sig
}

// Interrupt inplement the interface
func (c *CancelInterruptor) Interrupt() InterruptHandler {
	return func(ctx context.Context) bool {
		//always exit
		return true
	}
}

func init() {
	DefaultCancelInterruptors = []Interruptor{
		NewCancelInterruptor(syscall.SIGINT),
		NewCancelInterruptor(syscall.SIGTERM),
		NewCancelInterruptor(syscall.SIGKILL),
	}
}
