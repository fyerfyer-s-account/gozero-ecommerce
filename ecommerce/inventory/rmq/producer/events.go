package producer

import (
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
)

func NewStockUpdatedEvent(
    warehouseID int64,
    skuID int64,
    oldQuantity int32,
    newQuantity int32,
    reason string,
) *types.StockUpdatedEvent {
    return &types.StockUpdatedEvent{
        InventoryEvent: types.InventoryEvent{
            Type:        types.StockUpdated,
            WarehouseID: warehouseID,
            Timestamp:   time.Now(),
        },
        SkuID:       skuID,
        OldQuantity: oldQuantity,
        NewQuantity: newQuantity,
        Reason:      reason,
    }
}

func NewStockLockedEvent(
    warehouseID int64,
    orderNo string,
    items []types.StockItem,
) *types.StockLockedEvent {
    return &types.StockLockedEvent{
        InventoryEvent: types.InventoryEvent{
            Type:        types.StockLocked,
            WarehouseID: warehouseID,
            Timestamp:   time.Now(),
        },
        OrderNo: orderNo,
        Items:   items,
    }
}

func NewStockUnlockedEvent(
    warehouseID int64,
    orderNo string,
    items []types.StockItem,
) *types.StockUnlockedEvent {
    return &types.StockUnlockedEvent{
        InventoryEvent: types.InventoryEvent{
            Type:        types.StockUnlocked,
            WarehouseID: warehouseID,
            Timestamp:   time.Now(),
        },
        OrderNo: orderNo,
        Items:   items,
    }
}

func NewStockAlertEvent(
    warehouseID int64,
    skuID int64,
    current int32,
    threshold int32,
) *types.StockAlertEvent {
    return &types.StockAlertEvent{
        InventoryEvent: types.InventoryEvent{
            Type:        types.StockAlert,
            WarehouseID: warehouseID,
            Timestamp:   time.Now(),
        },
        SkuID:     skuID,
        Current:   current,
        Threshold: threshold,
    }
}

func NewStockOutOfStockEvent(
    warehouseID int64,
    skuID int64,
    quantity int32,
    reason string,
) *types.StockOutOfStockEvent {
    return &types.StockOutOfStockEvent{
        InventoryEvent: types.InventoryEvent{
            Type:        types.StockOutOfStock,
            WarehouseID: warehouseID,
            Timestamp:   time.Now(),
        },
        SkuID:    skuID,
        Quantity: quantity,
        Reason:   reason,
    }
}

func NewStockLowStockEvent(
    warehouseID int64,
    skuID int64,
    quantity int32,
    threshold int32,
) *types.StockLowStockEvent {
    return &types.StockLowStockEvent{
        InventoryEvent: types.InventoryEvent{
            Type:        types.StockLowStock,
            WarehouseID: warehouseID,
            Timestamp:   time.Now(),
        },
        SkuID:     skuID,
        Quantity:  quantity,
        Threshold: threshold,
    }
}