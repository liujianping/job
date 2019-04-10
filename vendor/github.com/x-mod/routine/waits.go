package routine

import (
	"context"
	"sync"
)

type _wait struct{}

// WithWait context with sync.WaitGroup
func WithWait(ctx context.Context) context.Context {
	if ctx != nil {
		return context.WithValue(ctx, _wait{}, &sync.WaitGroup{})
	}
	return context.WithValue(context.TODO(), _wait{}, &sync.WaitGroup{})
}

// WaitAdd if context with sync.WaitGroup, wait.Add
func WaitAdd(ctx context.Context, delta int) {
	if ctx != nil {
		wait := ctx.Value(_wait{})
		if wait != nil {
			wait.(*sync.WaitGroup).Add(delta)
		}
	}
}

// WaitDone if context with sync.WaitGroup, wait.Done
func WaitDone(ctx context.Context) {
	if ctx != nil {
		wait := ctx.Value(_wait{})
		if wait != nil {
			wait.(*sync.WaitGroup).Done()
		}
	}
}

// Wait should be invoked when Executor implemention use Go
func Wait(ctx context.Context) {
	if ctx != nil {
		wait := ctx.Value(_wait{})
		if wait != nil {
			wait.(*sync.WaitGroup).Wait()
		}
	}
}
