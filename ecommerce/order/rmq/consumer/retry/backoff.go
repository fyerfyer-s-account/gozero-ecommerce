package retry

import (
    "math"
    "math/rand"
    "time"
)

type BackoffStrategy interface {
    NextBackoff(attempt int) time.Duration
}

type ExponentialBackoff struct {
    initialInterval time.Duration
    maxInterval    time.Duration
    factor         float64
    jitter        bool
}

func NewExponentialBackoff(initialInterval, maxInterval time.Duration, factor float64, jitter bool) *ExponentialBackoff {
    return &ExponentialBackoff{
        initialInterval: initialInterval,
        maxInterval:    maxInterval,
        factor:         factor,
        jitter:        jitter,
    }
}

func (b *ExponentialBackoff) NextBackoff(attempt int) time.Duration {
    if attempt < 0 {
        attempt = 0
    }

    backoff := float64(b.initialInterval) * math.Pow(b.factor, float64(attempt))
    if backoff > float64(b.maxInterval) {
        backoff = float64(b.maxInterval)
    }

    if b.jitter {
        backoff = backoff * (1 + rand.Float64()*0.1 - 0.05) // ±5% jitter
    }

    return time.Duration(backoff)
}