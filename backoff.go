package retry

import (
	"context"
	"time"
)

type backoffStrategy struct {
	duration time.Duration
	backoff  float64
}

func (strategy *backoffStrategy) Next(ctx context.Context) bool {
	timer := time.NewTimer(strategy.duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		strategy.duration = time.Duration(float64(strategy.duration) * strategy.backoff)
		return true
	}
}

// WithBackoff .
func WithBackoff(duration time.Duration, backoff float64) Option {
	if backoff < 1 {
		panic("backoff param must > 1")
	}

	return func(action *actionImpl) {
		action.backoff = &backoffStrategy{
			duration: duration,
			backoff:  backoff,
		}
	}
}
