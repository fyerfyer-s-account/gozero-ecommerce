package types

import "errors"

var (
    // Event errors
    ErrEmptyEventID      = errors.New("empty event id")
    ErrEmptyEventType    = errors.New("empty event type")
    ErrEmptyTimestamp    = errors.New("empty timestamp")
    ErrEmptyEventData    = errors.New("empty event data")
    ErrEmptyTraceID      = errors.New("empty trace id")
    
    // Queue errors
    ErrQueueNotFound     = errors.New("queue not found")
    ErrQueueClosed       = errors.New("queue closed")
    
    // Processing errors
    ErrMaxRetriesExceeded = errors.New("max retries exceeded")
    ErrProcessingTimeout  = errors.New("processing timeout")
    ErrInvalidMessage     = errors.New("invalid message")
)

// RetryableError wraps an error that can be retried
type RetryableError struct {
    Err error
}

func (e *RetryableError) Error() string {
    return e.Err.Error()
}

func NewRetryableError(err error) *RetryableError {
    return &RetryableError{Err: err}
}

func IsRetryable(err error) bool {
    if err == nil {
        return false
    }
    _, ok := err.(*RetryableError)
    return ok
}