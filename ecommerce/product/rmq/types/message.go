package types

import "time"

// EventType defines the type of product event
type EventType string

const (
	// Product events
	EventTypeProductCreated       EventType = "product.created"
	EventTypeProductUpdated       EventType = "product.updated"
	EventTypeProductDeleted       EventType = "product.deleted"
	EventTypeProductStatusChanged EventType = "product.status.changed"

	// Price events
	EventTypePriceUpdated   EventType = "price.updated"
	EventTypePriceDropped   EventType = "price.dropped"
	EventTypePriceIncreased EventType = "price.increased"

	// Stock events
	EventTypeStockUpdated EventType = "stock.updated"
	EventTypeStockLow     EventType = "stock.low"
	EventTypeStockOut     EventType = "stock.out"
	EventTypeStockIn      EventType = "stock.in"

	// SKU events
	EventTypeSkuCreated      EventType = "sku.created"
	EventTypeSkuUpdated      EventType = "sku.updated"
	EventTypeSkuStockUpdated EventType = "sku.stock.updated"
	EventTypeSkuPriceUpdated EventType = "sku.price.updated"

	// Review events
	EventTypeReviewCreated       EventType = "review.created"
	EventTypeReviewUpdated       EventType = "review.updated"
	EventTypeReviewStatusChanged EventType = "review.status.changed"
	EventTypeReviewDeleted       EventType = "review.deleted"
)

// ProductEvent represents a product-related event message
type ProductEvent struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	ProductID int64     `json:"productId"`
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

// ProductData represents data for product events
type ProductData struct {
	Name        string   `json:"name,omitempty"`
	Brief       string   `json:"brief,omitempty"`
	Description string   `json:"description,omitempty"`
	CategoryId  int64    `json:"categoryId,omitempty"`
	Brand       string   `json:"brand,omitempty"`
	Images      []string `json:"images,omitempty"`
	Price       float64  `json:"price,omitempty"`
	Stock       int32    `json:"stock,omitempty"`
	Status      int32    `json:"status,omitempty"`
}

// SkuData represents data for SKU events
type SkuData struct {
	ID         int64             `json:"id"`
	ProductID  int64             `json:"productId"`
	Code       string            `json:"code"`
	Price      float64           `json:"price"`
	Stock      int32             `json:"stock"`
	Attributes map[string]string `json:"attributes"`
}

// ReviewData represents data for review events
type ReviewData struct {
	ID        int64    `json:"id"`
	ProductID int64    `json:"productId"`
	OrderID   int64    `json:"orderId"`
	UserID    int64    `json:"userId"`
	Rating    int32    `json:"rating"`
	Content   string   `json:"content"`
	Images    []string `json:"images,omitempty"`
	Status    int32    `json:"status,omitempty"`
}

// PriceData represents data for price events
type PriceData struct {
	ProductID    int64   `json:"productId"`
	OldPrice     float64 `json:"oldPrice"`
	NewPrice     float64 `json:"newPrice"`
	ChangeRatio  float64 `json:"changeRatio"`
	EffectiveAt  int64   `json:"effectiveAt"`
	CurrencyCode string  `json:"currencyCode"`
}

// StockData represents data for stock events
type StockData struct {
	ProductID   int64 `json:"productId"`
	OldStock    int32 `json:"oldStock"`
	NewStock    int32 `json:"newStock"`
	Threshold   int32 `json:"threshold,omitempty"`
	UpdatedAt   int64 `json:"updatedAt"`
	WarehouseID int64 `json:"warehouseId,omitempty"`
}

// NewProductEvent creates a new product event
func NewProductEvent(eventType EventType, productID int64, data any) *ProductEvent {
	return &ProductEvent{
		ID:        GenerateEventID(),
		Type:      eventType,
		ProductID: productID,
		Timestamp: time.Now().UnixMilli(),
		Data:      data,
		Metadata: Metadata{
			Source:  "product.service",
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
