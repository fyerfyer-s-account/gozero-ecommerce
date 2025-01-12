package handlers

import (
	"context"
	"fmt"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventoryclient"
)

type StockLockHandler struct {
    inventoryRpc inventoryclient.Inventory
}

func NewStockLockHandler(inventoryRpc inventoryclient.Inventory) *StockLockHandler {
    return &StockLockHandler{
        inventoryRpc: inventoryRpc,
    }
}

func (h *StockLockHandler) Handle(event *types.InventoryEvent) error {
    data, ok := event.Data.(*types.StockLockData)
    if !ok {
        return fmt.Errorf("invalid event data type")
    }

    switch event.Type {
    case types.EventTypeStockLocked:
        items := make([]*inventoryclient.LockItem, len(data.Items))
        for i, item := range data.Items {
            items[i] = &inventoryclient.LockItem{
                SkuId:       int64(item.SkuID),
                WarehouseId: int64(item.WarehouseID),
                Quantity:    item.Quantity,
            }
        }
        
        _, err := h.inventoryRpc.LockStock(context.Background(), &inventoryclient.LockStockRequest{
            OrderNo: data.OrderNo,
            Items:  items,
        })
        return err

    case types.EventTypeStockUnlocked:
        _, err := h.inventoryRpc.UnlockStock(context.Background(), &inventoryclient.UnlockStockRequest{
            OrderNo: data.OrderNo,
        })
        return err

    default:
        return fmt.Errorf("unknown event type: %s", event.Type)
    }
}