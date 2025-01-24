package types

import "time"

// Metadata represents event metadata for tracing and processing
type Metadata struct {
	TraceID      string            `json:"trace_id"`      // Unique trace ID for the event
	Source       string            `json:"source"`        // Source service name
	Timestamp    time.Time         `json:"timestamp"`     // Event creation time
	Version      string            `json:"version"`       // Event version
	Headers      map[string]string `json:"headers"`       // Custom headers
	RetryCount   int               `json:"retry_count"`   // Number of retries
	DeliveryMode string            `json:"delivery_mode"` // persistent or non-persistent
}

// MetadataBuilder helps build metadata with a fluent interface
type MetadataBuilder struct {
	metadata *Metadata
}

// NewMetadataBuilder creates a new MetadataBuilder
func NewMetadataBuilder() *MetadataBuilder {
	return &MetadataBuilder{
		metadata: &Metadata{
			Timestamp: time.Now(),
			Headers:   make(map[string]string),
		},
	}
}
