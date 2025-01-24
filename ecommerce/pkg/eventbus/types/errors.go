package types

// EventError represents a base error type for event processing
type EventError struct {
    Code    string
    Message string
    Err     error
}

func (e *EventError) Error() string {
    if e.Err != nil {
        return e.Message + ": " + e.Err.Error()
    }
    return e.Message
}

// RetryableError indicates the error can be retried
type RetryableError struct {
    *EventError
}

// NonRetryableError indicates the error cannot be retried
type NonRetryableError struct {
    *EventError
}