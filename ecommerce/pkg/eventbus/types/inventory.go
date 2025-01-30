package types

import "time"

type InventoryEventType string

const (
	StockLocked             InventoryEventType = "inventory.stock.locked"
	StockUnlocked           InventoryEventType = "inventory.stock.unlocked"
	StockDeducted           InventoryEventType = "inventory.stock.deducted"
	StockIncremented        InventoryEventType = "inventory.stock.incremented"
	StockAlert              InventoryEventType = "inventory.stock.alert"
	StockChecked            InventoryEventType = "inventory.stock.checked"
	StockReserved           InventoryEventType = "inventory.stock.reserved"
	StockUpdated            InventoryEventType = "inventory.stock.updated"
	StockOutOfStock         InventoryEventType = "inventory.stock.out_of_stock"
	StockLowStock           InventoryEventType = "inventory.stock.low_stock"
	InventoryPaymentSuccess InventoryEventType = "inventory.payment.success"
	InventoryPaymentFailed  InventoryEventType = "inventory.payment.failed"
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

// StockCheckedEvent represents stock check result event
type StockCheckedEvent struct {
	InventoryEvent
	OrderNo string      `json:"order_no"`
	Items   []StockItem `json:"items"`
	Success bool        `json:"success"`
	Reason  string      `json:"reason,omitempty"`
}

// StockReservedEvent represents stock reservation event
type StockReservedEvent struct {
	InventoryEvent
	OrderNo    string      `json:"order_no"`
	Items      []StockItem `json:"items"`
	ExpireTime time.Time   `json:"expire_time"`
}

// Add validation methods
func (e *StockLockedEvent) Validate() error {
	if e.OrderNo == "" || len(e.Items) == 0 {
		return &NonRetryableError{
			EventError: &EventError{
				Code:    "INVALID_STOCK_LOCK_EVENT",
				Message: "order_no and items are required",
			},
		}
	}
	return nil
}

type StockUpdatedEvent struct {
	InventoryEvent
	SkuID       int64  `json:"sku_id"`
	OldQuantity int32  `json:"old_quantity"`
	NewQuantity int32  `json:"new_quantity"`
	Reason      string `json:"reason"`
}

type StockOutOfStockEvent struct {
	InventoryEvent
	SkuID    int64  `json:"sku_id"`
	Quantity int32  `json:"quantity"`
	Reason   string `json:"reason"`
}

type StockLowStockEvent struct {
	InventoryEvent
	SkuID     int64 `json:"sku_id"`
	Quantity  int32 `json:"quantity"`
	Threshold int32 `json:"threshold"`
}

// Add validation methods for new events
func (e *StockUpdatedEvent) Validate() error {
	if e.SkuID == 0 {
		return &NonRetryableError{
			EventError: &EventError{
				Code:    "INVALID_STOCK_UPDATE_EVENT",
				Message: "sku_id is required",
			},
		}
	}
	return nil
}

// InventoryPaymentSuccessEvent represents inventory update after successful payment
type InventoryPaymentSuccessEvent struct {
    InventoryEvent
    OrderNo      string      `json:"order_no"`
    PaymentNo    string      `json:"payment_no"`
    Items        []StockItem `json:"items"`
    DeductTime   time.Time   `json:"deduct_time"`
}

// InventoryPaymentFailedEvent represents inventory release after payment failure
type InventoryPaymentFailedEvent struct {
    InventoryEvent
    OrderNo      string      `json:"order_no"`
    PaymentNo    string      `json:"payment_no"`
    Items        []StockItem `json:"items"`
    ReleaseTime  time.Time   `json:"release_time"`
    Reason       string      `json:"reason"`
}

// Add validation methods
func (e *InventoryPaymentSuccessEvent) Validate() error {
    if e.OrderNo == "" || e.PaymentNo == "" || len(e.Items) == 0 {
        return &NonRetryableError{
            EventError: &EventError{
                Code:    "INVALID_INVENTORY_PAYMENT_SUCCESS",
                Message: "order_no, payment_no and items are required",
            },
        }
    }
    return nil
}

func (e *InventoryPaymentFailedEvent) Validate() error {
    if e.OrderNo == "" || e.PaymentNo == "" || len(e.Items) == 0 {
        return &NonRetryableError{
            EventError: &EventError{
                Code:    "INVALID_INVENTORY_PAYMENT_FAILED",
                Message: "order_no, payment_no and items are required",
            },
        }
    }
    return nil
}