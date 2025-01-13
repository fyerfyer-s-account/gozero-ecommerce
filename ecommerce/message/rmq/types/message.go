package types

import (
	"time"
)

const (
	// Message status
	MessageStatusPending  = 1
	MessageStatusSending  = 2
	MessageStatusSuccess  = 3 
	MessageStatusFailed   = 4

	// Message types
	MessageTypeSystem   = 1
	MessageTypeOrder    = 2 
	MessageTypeActivity = 3
	MessageTypeLogistics = 4

	// Send channels
	ChannelInApp  = 1
	ChannelSMS    = 2
	ChannelEmail  = 3
	ChannelPush   = 4
)

type EventType string

const (
	EventTypeMessageCreated  EventType = "message.created"
	EventTypeMessageSent     EventType = "message.sent"
	EventTypeMessageRead     EventType = "message.read"
	EventTypeTemplateCreated EventType = "template.created"
	EventTypeTemplateUpdated EventType = "template.updated"
)

type MessageEvent struct {
	ID        string      `json:"id"`
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
	Metadata  Metadata    `json:"metadata"`
}

type Metadata struct {
	Source  string            `json:"source"`
	UserID  int64             `json:"userId,omitempty"`
	TraceID string            `json:"traceId"`
	Tags    map[string]string `json:"tags,omitempty"`
}

type MessageCreatedData struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userId"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Type        int32  `json:"type"`
	SendChannel int32  `json:"sendChannel"`
	ExtraData   string `json:"extraData,omitempty"`
}

type MessageSentData struct {
	MessageID   int64  `json:"messageId"`
	UserID      int64  `json:"userId"` 
	SendChannel int32  `json:"sendChannel"`
	Status      int32  `json:"status"` // Changed to int32
	Error       string `json:"error,omitempty"`
	RetryCount  int32  `json:"retryCount"`
}

type MessageReadData struct {
	MessageID int64     `json:"messageId"`
	UserID    int64     `json:"userId"`
	ReadTime  time.Time `json:"readTime"`
}

type TemplateData struct {
	ID              int64   `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	TitleTemplate   string  `json:"titleTemplate"`
	ContentTemplate string  `json:"contentTemplate"`
	Type            int32   `json:"type"`
	Channels        []int32 `json:"channels"`
	Config          string  `json:"config"`
}
