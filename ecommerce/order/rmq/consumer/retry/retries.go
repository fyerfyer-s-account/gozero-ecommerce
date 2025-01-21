package retry

import (
    "context"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/types"
    "time"
)

type Retrier struct {
    maxAttempts     int
    backoff         BackoffStrategy
    shouldRetry     func(error) bool
}

func NewRetrier(maxAttempts int, backoff BackoffStrategy) *Retrier {
    return &Retrier{
        maxAttempts: maxAttempts,
        backoff:     backoff,
        shouldRetry: types.IsRetryable,
    }
}

func (r *Retrier) Do(operation func() error) error {
    var lastErr error
    
    for attempt := 0; attempt < r.maxAttempts; attempt++ {
        err := operation()
        if err == nil {
            return nil
        }

        lastErr = err
        if !r.shouldRetry(err) {
            return err
        }

        if attempt == r.maxAttempts-1 {
            break
        }

        backoffDuration := r.backoff.NextBackoff(attempt)
        timer := time.NewTimer(backoffDuration)
        select {
        case <-timer.C:
            continue
        }
        timer.Stop()
    }

    if lastErr != nil {
        return &types.RetryableError{
            Err: types.ErrMaxRetriesExceeded,
        }
    }

    return nil
}

func (r *Retrier) DoWithContext(ctx context.Context, operation func() error) error {
    var lastErr error
    
    for attempt := 0; attempt < r.maxAttempts; attempt++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            err := operation()
            if err == nil {
                return nil
            }

            lastErr = err
            if !r.shouldRetry(err) {
                return err
            }

            if attempt == r.maxAttempts-1 {
                break
            }

            backoffDuration := r.backoff.NextBackoff(attempt)
            timer := time.NewTimer(backoffDuration)
            select {
            case <-timer.C:
                continue
            case <-ctx.Done():
                timer.Stop()
                return ctx.Err()
            }
        }
    }

    if lastErr != nil {
        return &types.RetryableError{
            Err: types.ErrMaxRetriesExceeded,
        }
    }

    return nil
}