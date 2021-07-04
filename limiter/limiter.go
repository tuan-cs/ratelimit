package limiter

import (
	"context"

	"github.com/ratelimit/runner"
)

type Config struct {
	IdentifierExtrator func(ctx context.Context) (string, error)

	ErrorHandler func(ctx context.Context, err error) error

	DenyHandler func(ctx context.Context, identifier string, err error) error

	Store LimiterStore
}

type LimiterStore interface {
	// Stores for the rate limiter have to implement the Allow method
	Allow(identifier string) (bool, error)
}

// New funtion returns a new limiter ready executor, the execution will be limited
// the number of times on the config
func New(config Config) runner.Runner {
	return NewMiddleware(config)(nil)
}

// NewMiddleware returns a new ratelimit middleware, the execution will be
// limited the number of times on the config
func NewMiddleware(config Config) runner.Middleware {

	return func(next runner.Runner) runner.Runner {
		next = runner.SanitizeRunner(next)

		return runner.RunnerFunc(func(ctx context.Context, f runner.Func) error {
			var err error
			identifier, err := config.IdentifierExtrator(ctx)
			if err != nil {
				config.ErrorHandler(ctx, err)
				return err
			}

			if allow, err := config.Store.Allow(identifier); !allow {
				config.DenyHandler(ctx, identifier, err)
				return err
			}

			err = next.Run(ctx, f)
			if err != nil {
				return err
			}

			return err
		})
	}
}
