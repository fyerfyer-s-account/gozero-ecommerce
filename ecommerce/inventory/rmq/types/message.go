package types

import (
	"fmt"
	"math/rand"
	"time"
)

// EventType defines the type of inventory event
type EventType string

const (
	EventTypeStockUpdated   EventType = "stock.updated"
	EventTypeStockAlert     EventType = "stock.alert"
	EventTypeStockLocked    EventType = "stock.locked"
	EventTypeStockUnlocked  EventType = "stock.unlocked"
	EventTypeOrderCreated   EventType = "order.created"
	EventTypeOrderPaid      EventType = "order.paid"
	EventTypeOrderCancelled EventType = "order.cancelled"
	EventTypeOrderRefunded  EventType = "order.refund.completed"
)

// InventoryEvent represents an inventory-related event message
type InventoryEvent struct {
	ID        string      `json:"id"`
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
	Metadata  Metadata    `json:"metadata"`
}

// Metadata contains additional information about the event
type Metadata struct {
	Source  string            `json:"source"`
	UserID  int64             `json:"userId,omitempty"`
	TraceID string            `json:"traceId"`
	Tags    map[string]string `json:"tags,omitempty"`
}

// StockAlertData represents data for stock alert events
type StockAlertData struct {
	SkuID       uint64 `json:"skuId"`
	WarehouseID uint64 `json:"warehouseId"`
	Available   int64  `json:"available"`
	Threshold   int64  `json:"threshold"`
}

type StockLockData struct {
	OrderNo string     `json:"orderNo"`
	Items   []LockItem `json:"items"`
}

type LockItem struct {
	SkuID       uint64 `json:"skuId"`
	WarehouseID uint64 `json:"warehouseId"`
	Quantity    int32  `json:"quantity"`
}

// StockUpdateData represents data for stock update events
type StockUpdateData struct {
	SkuID       uint64 `json:"skuId"`
	WarehouseID uint64 `json:"warehouseId"`
	Quantity    int32  `json:"quantity"`
	Remark      string `json:"remark"`
}

type OrderCreatedData struct {
	OrderNo string      `json:"orderNo"`
	Items   []OrderItem `json:"items"`
}

type OrderItem struct {
	SkuID       uint64 `json:"skuId"`
	WarehouseID uint64 `json:"warehouseId"`
	Quantity    int32  `json:"quantity"`
}

type OrderCancelledData struct {
	OrderNo string `json:"orderNo"`
}

type OrderPaidData struct {
	OrderNo string `json:"orderNo"`
}

type OrderRefundedData struct {
	OrderNo string      `json:"orderNo"`
	Items   []OrderItem `json:"items"`
}

// NewInventoryEvent creates a new inventory event
func NewInventoryEvent(eventType EventType, data interface{}, userId int64) *InventoryEvent {
	return &InventoryEvent{
		ID:        GenerateEventID(),
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
		Metadata: Metadata{
			Source:  "inventory-service",
			UserID:  userId,
			TraceID: GenerateEventID(),
			Tags:    make(map[string]string),
		},
	}
}

// GenerateEventID generates a unique event ID
func GenerateEventID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), RandomString(6))
}

// RandomString generates a random string of given length
func RandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
