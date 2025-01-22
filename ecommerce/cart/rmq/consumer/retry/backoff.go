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
    backoff := float64(b.initialInterval) * math.Pow(b.factor, float64(attempt))
    
    if b.jitter {
        backoff = backoff * (1 + rand.Float64()*0.1) // 10% jitter
    }
    
    if backoff > float64(b.maxInterval) {
        backoff = float64(b.maxInterval)
    }
    
    return time.Duration(backoff)
}