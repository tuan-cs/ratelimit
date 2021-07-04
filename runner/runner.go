package runner

import (
	"context"

	"github.com/ratelimit/errors"
)

// Func is the function to be execute
type Func func(c context.Context) error

// command is the unit of execution
type command struct{}

func (c command) Run(ctx context.Context, f Func) error {
	// Only execute if we reached to the execution and the context has not been canceled
	select {
	case <-ctx.Done():
		return errors.ErrContextCanceled
	default:
		return f(ctx)
	}
}

type Runner interface {
	// Run will run the unit of execution passed on f
	Run(ctx context.Context, f Func) error
}

// Runnerfunc is a helper that will satisfies limit.Limiter interface by using a function
type RunnerFunc func(ctx context.Context, f Func) error

// Run satisfies Runner interface
func (r RunnerFunc) Run(ctx context.Context, f Func) error {
	return r(ctx, f)
}

// Middleware represents a middleware for a runner, it takes a runner and returns a runner
type Middleware func(Runner) Runner

// RunnerChain will get N middleware for a runner, and will create a Runner
// chain with them in the order that have been passed
func RunnerChain(middlewares ...Middleware) Runner {
	var runner Runner = &command{}

	// Start wrapping in reverse order
	for i := len(middlewares) - 1; i >= 0; i-- {
		runner = middlewares[i](runner)
	}

	return runner
}

// SanitizeRunner returns a safe execution Runner if the runner is nil. Usually
// this helper will be used for the last part of the runner chain when the
// runner is nil, so instead of acting on a nil Runner its executed on a
// `command` Runner, this runner knows how to execute the `Func` function.
// It's safe to use it always as if it encounters a safe Runner it will return
// that Runner
func SanitizeRunner(r Runner) Runner {
	if r == nil {
		return &command{}
	}
	return r
}
