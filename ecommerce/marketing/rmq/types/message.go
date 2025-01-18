package types

import "time"

// EventType defines the type of marketing event
type EventType string

const (
    // Coupon events
    EventTypeCouponCreated  EventType = "coupon.created"
    EventTypeCouponReceived EventType = "coupon.received"
    EventTypeCouponUsed     EventType = "coupon.used"
    EventTypeCouponExpired  EventType = "coupon.expired"

    // Promotion events
    EventTypePromotionCreated EventType = "promotion.created"
    EventTypePromotionStarted EventType = "promotion.started"
    EventTypePromotionEnded   EventType = "promotion.ended"

    // Points events
    EventTypePointsAdded   EventType = "points.added"
    EventTypePointsUsed    EventType = "points.used"
    EventTypePointsExpired EventType = "points.expired"
)

// MarketingEvent represents a marketing-related event message
type MarketingEvent struct {
    ID        string    `json:"id"`
    Type      EventType `json:"type"`
    Timestamp int64     `json:"timestamp"`
    Data      any       `json:"data"`
    Metadata  Metadata  `json:"metadata"`
}

// Metadata contains additional information about the event
type Metadata struct {
    UserID  int64  `json:"userId,omitempty"`
    TraceID string `json:"traceId,omitempty"`
    Source  string `json:"source"`
    Version string `json:"version"`
}

// CouponEventData represents data for coupon events
type CouponEventData struct {
    CouponID    int64   `json:"couponId"`
    Code        string  `json:"code"`
    Type        int32   `json:"type"`
    Value       float64 `json:"value"`
    MinAmount   float64 `json:"minAmount"`
    Status      int32   `json:"status"`
    UserID      int64   `json:"userId,omitempty"`
    OrderNo     string  `json:"orderNo,omitempty"`
}

// PromotionEventData represents data for promotion events
type PromotionEventData struct {
    PromotionID int64  `json:"promotionId"`
    Name        string `json:"name"`
    Type        int32  `json:"type"`
    Rules       string `json:"rules"`
    Status      int32  `json:"status"`
    StartTime   int64  `json:"startTime"`
    EndTime     int64  `json:"endTime"`
}

// PointsEventData represents data for points events
type PointsEventData struct {
    UserID      int64  `json:"userId"`
    Points      int64  `json:"points"`
    Type        int32  `json:"type"`
    Source      string `json:"source"`
    Remark      string `json:"remark"`
    OrderNo     string `json:"orderNo,omitempty"`
}

// NewMarketingEvent creates a new marketing event
func NewMarketingEvent(eventType EventType, data any) *MarketingEvent {
    return &MarketingEvent{
        ID:        GenerateEventID(),
        Type:      eventType,
        Timestamp: time.Now().UnixMilli(),
        Data:      data,
        Metadata: Metadata{
            Source:  "marketing.service",
            Version: "1.0",
        },
    }
}

// GenerateEventID generates a unique event ID
func GenerateEventID() string {
    return time.Now().Format("20060102150405.000") + "-" + RandomString(8)
}

// RandomString generates a random string of given length
func RandomString(n int) string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
    }
    return string(b)
}