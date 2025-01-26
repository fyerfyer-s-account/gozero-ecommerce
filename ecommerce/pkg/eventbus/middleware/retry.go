package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/streadway/amqp"
)

type RetryConfig struct {
    MaxRetries  int
    InitialWait time.Duration
    MaxWait     time.Duration
}

func DefaultRetryConfig() *RetryConfig {
    return &RetryConfig{
        MaxRetries:  3,
        InitialWait: time.Second,
        MaxWait:     time.Second * 30,
    }
}

func Retry(config *RetryConfig) MiddlewareFunc {
    if config == nil {
        config = DefaultRetryConfig()
    }

    return func(next HandlerFunc) HandlerFunc {
        return func(ctx context.Context, msg amqp.Delivery) error {
            attempts := 0
            var lastErr error

            for attempts <= config.MaxRetries {
                if err := next(ctx, msg); err != nil {
                    if _, ok := err.(*types.NonRetryableError); ok {
                        return err // Don't retry non-retryable errors
                    }

                    lastErr = err
                    attempts++

                    if attempts <= config.MaxRetries {
                        wait := exponentialBackoff(attempts, config.InitialWait, config.MaxWait)
                        time.Sleep(wait)
                        continue
                    }
                }
                return nil // Success
            }

            return fmt.Errorf("max retries reached: %v", lastErr)
        }
    }
}

func exponentialBackoff(attempt int, initial, max time.Duration) time.Duration {
    wait := initial * time.Duration(1<<uint(attempt-1))
    if wait > max {
        return max
    }
    return wait
}