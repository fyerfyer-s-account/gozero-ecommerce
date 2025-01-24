package types

import "time"

type InventoryEventType string

const (
	StockLocked      InventoryEventType = "inventory.stock.locked"
	StockUnlocked    InventoryEventType = "inventory.stock.unlocked"
	StockDeducted    InventoryEventType = "inventory.stock.deducted"
	StockIncremented InventoryEventType = "inventory.stock.incremented"
	StockAlert       InventoryEventType = "inventory.stock.alert"
)

// InventoryEvent represents the base inventory event structure
type InventoryEvent struct {
	Type        InventoryEventType `json:"type"`
	WarehouseID int64              `json:"warehouse_id"`
	Timestamp   time.Time          `json:"timestamp"`
}

// StockItem represents a stock item in inventory events
type StockItem struct {
	SkuID    int64 `json:"sku_id"`
	Quantity int32 `json:"quantity"`
}

// StockLockedEvent represents stock locking event
type StockLockedEvent struct {
	InventoryEvent
	OrderNo string      `json:"order_no"`
	Items   []StockItem `json:"items"`
}

// StockUnlockedEvent represents stock unlocking event
type StockUnlockedEvent struct {
	InventoryEvent
	OrderNo string      `json:"order_no"`
	Items   []StockItem `json:"items"`
}

// StockDeductedEvent represents stock deduction event
type StockDeductedEvent struct {
	InventoryEvent
	OrderNo string      `json:"order_no"`
	Items   []StockItem `json:"items"`
}

// StockAlertEvent represents stock alert event
type StockAlertEvent struct {
	InventoryEvent
	SkuID     int64 `json:"sku_id"`
	Current   int32 `json:"current"`
	Threshold int32 `json:"threshold"`
}
