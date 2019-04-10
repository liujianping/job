package routine

import (
	"context"

	"github.com/x-mod/errors"
)

var (
	//ErrNoneExecutor error
	ErrNoneExecutor = errors.New("none executor")
	//ErrNoneContext error
	ErrNoneContext = errors.New("none context")
	//ErrNonePlan error
	ErrNonePlan = errors.New("none plan")
)

//Executor interface definition
type Executor interface {
	//Execute before stopping make sure all subroutines stopped
	Execute(context.Context) error
}

//ExecutorFunc definition
type ExecutorFunc func(context.Context) error

//Execute ExecutorFunc implemention of Executor
func (f ExecutorFunc) Execute(ctx context.Context) error {
	return f(ctx)
}

// ExecutorMiddleware is a function that middlewares can implement to be
// able to chain.
type ExecutorMiddleware func(Executor) Executor

// UseExecutorMiddleware wraps a Executor in one or more middleware.
func UseExecutorMiddleware(exec Executor, middleware ...ExecutorMiddleware) Executor {
	// Apply in reverse order.
	for i := len(middleware) - 1; i >= 0; i-- {
		m := middleware[i]
		exec = m(exec)
	}
	return exec
}

//Routine for Executors
type Routine interface {
	Go(context.Context, Executor) chan error
}

//GoFunc definition
type GoFunc func(context.Context, Executor) chan error

//Go RunnerFunc implemention of Runner
func (f GoFunc) Go(ctx context.Context, exec Executor) chan error {
	return f(ctx, exec)
}
