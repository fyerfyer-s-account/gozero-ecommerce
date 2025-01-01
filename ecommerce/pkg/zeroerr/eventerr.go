package zeroerr

import "errors"

var (
	ErrInvalidEventData      = errors.New("invalid event data")
	ErrInvalidEventType      = errors.New("invalid event type")
	ErrEventHandlerNotFound  = errors.New("event handler not found")
	ErrEventProcessingFailed = errors.New("event processing failed")
)
