package routine

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/x-mod/errors"
)

//Main wrapper for executor with waits & signal interuptors
func Main(parent context.Context, exec Executor, opts ...Opt) error {
	moptions := &options{}
	for _, opt := range opts {
		opt(moptions)
	}
	// context with cancel & wait
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	// argments
	if len(moptions.args) > 0 {
		ctx = WithArgments(ctx, moptions.args...)
	}
	// signals
	sigchan := make(chan os.Signal)
	sighandlers := make(map[os.Signal]InterruptHandler)
	for _, interruptor := range moptions.interrupts {
		signal.Notify(sigchan, interruptor.Signal())
		sighandlers[interruptor.Signal()] = interruptor.Interrupt()
	}
	// main executor
	ch := Go(ctx, exec)
	// main exit for sig & finished
	exitCh := make(chan error, 1)
	for {
		select {
		case sig := <-sigchan:
			// cancel when a signal catched
			if h, ok := sighandlers[sig]; ok {
				if h(ctx) {
					exitCh <- errors.CodeError(SignalCode(sig.(syscall.Signal)))
					goto Exit
				}
			}
		case <-ctx.Done():
			exitCh <- errors.WithCode(ctx.Err(), GeneralErr)
			goto Exit
		case err := <-ch:
			exitCh <- err
			goto Exit
		}
	}
Exit:
	//exit hook
	if moptions.beforeExit != nil {
		moptions.beforeExit.Execute(ctx)
	}
	return <-exitCh
}

//Go wrapper for go keyword, use in MAIN function
func Go(ctx context.Context, exec Executor) chan error {
	ch := make(chan error, 1)
	if exec == nil {
		ch <- ErrNoneExecutor
		return ch
	}
	if ctx == nil {
		ch <- ErrNoneContext
		return ch
	}
	WaitAdd(ctx, 1)
	go func() {
		defer WaitDone(ctx)

		// channel for function (run) done
		stop := make(chan struct{})
		go func() {
			ch <- exec.Execute(ctx)
			close(stop)
		}()
		// run exit for cancel & finished
		select {
		case <-ctx.Done():
			ch <- ctx.Err()
		case <-stop:
		}
	}()
	return ch
}

type options struct {
	args          []interface{}
	interrupts    []Interruptor
	beforeExecute Executor
	afterExecute  Executor
	beforeExit    Executor
}

//Opt interface
type Opt func(*options)

//Arguments Opt for Main
func Arguments(args ...interface{}) Opt {
	return func(opts *options) {
		opts.args = args
	}
}

//Interrupt Opt for Main
func Interrupts(ints ...Interruptor) Opt {
	return func(opts *options) {
		opts.interrupts = append(opts.interrupts, ints...)
	}
}

//BeforeExit Opt for Main
func BeforeExit(exec Executor) Opt {
	return func(opts *options) {
		opts.beforeExit = exec
	}
}
