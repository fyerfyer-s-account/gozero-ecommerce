package retry

import (
	"context"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/types"
	"time"
)

type Retrier struct {
	maxAttempts int
	backoff     BackoffStrategy
	shouldRetry func(error) bool
}

func NewRetrier(maxAttempts int, backoff BackoffStrategy) *Retrier {
	return &Retrier{
		maxAttempts: maxAttempts,
		backoff:     backoff,
		shouldRetry: types.IsRetryable,
	}
}

func (r *Retrier) Do(operation func() error) error {

	for attempt := 0; attempt < r.maxAttempts; attempt++ {
		if err := operation(); err == nil {
			return nil
		} else {
			if !r.shouldRetry(err) {
				return err
			}
			if attempt < r.maxAttempts-1 {
				time.Sleep(r.backoff.NextBackoff(attempt))
			}
		}
	}

	return types.ErrMaxRetriesExceeded
}

func (r *Retrier) DoWithContext(ctx context.Context, operation func() error) error {
	for attempt := 0; attempt < r.maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := operation(); err == nil {
				return nil
			} else {
				if !r.shouldRetry(err) {
					return err
				}
				if attempt < r.maxAttempts-1 {
					timer := time.NewTimer(r.backoff.NextBackoff(attempt))
					select {
					case <-ctx.Done():
						timer.Stop()
						return ctx.Err()
					case <-timer.C:
					}
				}
			}
		}
	}

	return types.ErrMaxRetriesExceeded
}
