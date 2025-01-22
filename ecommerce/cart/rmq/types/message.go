package types

import (
    "time"
    "strconv"
)

type EventType string

const (
    EventTypeItemAdded    EventType = "cart.item.added"
    EventTypeItemUpdated  EventType = "cart.item.updated"
    EventTypeItemRemoved  EventType = "cart.item.removed"
    EventTypeCartCleared  EventType = "cart.cleared"
)

type CartEvent struct {
    ID        string        `json:"id"`
    Type      EventType     `json:"type"`
    Timestamp time.Time     `json:"timestamp"`
    Data      interface{}   `json:"data"`
    Metadata  EventMetadata `json:"metadata"`
}

type EventMetadata struct {
    TraceID   string            `json:"trace_id"`
    UserID    string            `json:"user_id"`
    Tags      map[string]string `json:"tags,omitempty"`
}

type CartItemData struct {
    UserID    int64   `json:"user_id"`
    ProductID int64   `json:"product_id"`
    Quantity  int32   `json:"quantity"`
    Selected  bool    `json:"selected"`
}

type CartItemAddedData struct {
    CartItemData
}

type CartItemUpdatedData struct {
    CartItemData
    OldQuantity int32 `json:"old_quantity"`
}

type CartItemRemovedData struct {
    UserID    int64 `json:"user_id"`
    ProductID int64 `json:"product_id"`
}

type CartClearedData struct {
    UserID int64 `json:"user_id"`
}

func (e *CartEvent) Validate() error {
    if e.ID == "" {
        return ErrEmptyEventID
    }
    if e.Type == "" {
        return ErrEmptyEventType
    }
    if e.Timestamp.IsZero() {
        return ErrEmptyTimestamp
    }
    if e.Data == nil {
        return ErrEmptyEventData
    }
    if e.Metadata.TraceID == "" {
        return ErrEmptyTraceID
    }
    if e.Metadata.UserID == "" {
        return ErrEmptyUserID
    }
    return nil
}

func (e *CartEvent) IsRetryable() bool {
    switch e.Type {
    case EventTypeItemAdded, EventTypeItemUpdated, EventTypeItemRemoved:
        return true
    default:
        return false
    }
}

func (e *CartEvent) GetRetryCount() int {
    if count, ok := e.Metadata.Tags["retry_count"]; ok {
        if v, err := strconv.Atoi(count); err == nil {
            return v
        }
    }
    return 0
}