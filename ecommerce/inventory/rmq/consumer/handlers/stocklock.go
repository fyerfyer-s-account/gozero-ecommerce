package handlers

import (
    "context"
    "fmt"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventoryclient"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/middleware"
)

type StockLockHandler struct {
    BaseHandler
    inventoryRpc inventoryclient.Inventory
}

func NewStockLockHandler(logger middleware.Logger, inventoryRpc inventoryclient.Inventory) *StockLockHandler {
    return &StockLockHandler{
        BaseHandler:   NewBaseHandler(logger),
        inventoryRpc: inventoryRpc,
    }
}

func (h *StockLockHandler) Handle(event *types.InventoryEvent) error {
    h.LogEvent(event, "handling stock lock event")
    
    data, ok := event.Data.(*types.StockLockData)
    if !ok {
        return fmt.Errorf("invalid event data type")
    }

    var err error
    switch event.Type {
    case types.EventTypeStockLocked:
        err = h.handleLock(context.Background(), data)
    case types.EventTypeStockUnlocked:
        err = h.handleUnlock(context.Background(), data)
    default:
        err = fmt.Errorf("unknown event type: %s", event.Type)
    }

    if err != nil {
        h.LogError(event, err)
        return types.NewRetryableError(err)
    }

    h.LogEvent(event, "stock lock event handled successfully")
    return nil
}

func (h *StockLockHandler) handleLock(ctx context.Context, data *types.StockLockData) error {
    items := make([]*inventoryclient.LockItem, len(data.Items))
    for i, item := range data.Items {
        items[i] = &inventoryclient.LockItem{
            SkuId:       int64(item.SkuID),
            WarehouseId: int64(item.WarehouseID),
            Quantity:    item.Quantity,
        }
    }
    
    _, err := h.inventoryRpc.LockStock(ctx, &inventoryclient.LockStockRequest{
        OrderNo: data.OrderNo,
        Items:  items,
    })
    return err
}

func (h *StockLockHandler) handleUnlock(ctx context.Context, data *types.StockLockData) error {
    _, err := h.inventoryRpc.UnlockStock(ctx, &inventoryclient.UnlockStockRequest{
        OrderNo: data.OrderNo,
    })
    return err
}