package types

// EventType represents the base type for all event types
type EventType string

// Event represents the base event structure
type Event struct {
    Type     EventType        `json:"type"`
    Data     interface{}      `json:"data"`
    Metadata *Metadata        `json:"metadata"`
}

// EventHandler represents a function that handles an event
type EventHandler func(event *Event) error

// EventProcessor defines the interface for processing events
type EventProcessor interface {
    Process(event *Event) error
}

// EventValidator defines the interface for validating events
type EventValidator interface {
    Validate() error
}

// BaseEvent provides common functionality for all events
type BaseEvent struct {
    Metadata *Metadata `json:"metadata"`
}

// WithMetadata adds metadata to an event
func (e *BaseEvent) WithMetadata(metadata *Metadata) {
    e.Metadata = metadata
}

// GetMetadata returns event metadata
func (e *BaseEvent) GetMetadata() *Metadata {
    return e.Metadata
}

// DeliveryMode defines message delivery guarantees
type DeliveryMode int

const (
    // NonPersistent messages may be lost
    NonPersistent DeliveryMode = 1
    // Persistent messages are written to disk
    Persistent DeliveryMode = 2
)

// Exchange types
const (
    DirectExchange  = "direct"
    FanoutExchange = "fanout"
    TopicExchange  = "topic"
    HeadersExchange = "headers"
)

// Common event status
const (
    EventStatusPending   = "pending"
    EventStatusSuccess  = "success"
    EventStatusFailed   = "failed"
    EventStatusRetrying = "retrying"
)