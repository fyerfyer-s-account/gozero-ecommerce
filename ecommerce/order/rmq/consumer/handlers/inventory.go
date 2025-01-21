package handlers

import (
	"context"
	"fmt"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventoryclient"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/middleware"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/types"
)

type InventoryHandler struct {
    BaseHandler
    inventoryRpc inventoryclient.Inventory
}

func NewInventoryHandler(logger middleware.Logger, inventoryRpc inventoryclient.Inventory) *InventoryHandler {
    return &InventoryHandler{
        BaseHandler:   NewBaseHandler(logger),
        inventoryRpc: inventoryRpc,
    }
}

func (h *InventoryHandler) Handle(event *types.OrderEvent) error {
    h.LogEvent(event, "handling inventory event")
    
    var err error
    switch event.Type {
    case types.EventTypeOrderCreated:
        err = h.handleOrderCreated(event)
    case types.EventTypeOrderCancelled:
        err = h.handleOrderCancelled(event)
    case types.EventTypeOrderPaid:
        err = h.handleOrderPaid(event)
    }

    if err != nil {
        h.LogError(event, err)
        return types.NewRetryableError(err)
    }
    return nil
}

func (h *InventoryHandler) handleOrderCreated(event *types.OrderEvent) error {
    data, ok := event.Data.(*types.OrderCreatedData)
    if !ok {
        return fmt.Errorf("invalid event data type")
    }

    items := make([]*inventoryclient.LockItem, 0)
    for _, item := range data.Items {
        items = append(items, &inventoryclient.LockItem{
            SkuId:    item.SkuID,
            Quantity: item.Quantity,
        })
    }

    _, err := h.inventoryRpc.LockStock(context.Background(), &inventoryclient.LockStockRequest{
        OrderNo: data.OrderNo,
        Items:   items,
    })
    return err
}

func (h *InventoryHandler) handleOrderCancelled(event *types.OrderEvent) error {
    data, ok := event.Data.(*types.OrderCancelledData)
    if !ok {
        return fmt.Errorf("invalid event data type")
    }

    _, err := h.inventoryRpc.UnlockStock(context.Background(), &inventoryclient.UnlockStockRequest{
        OrderNo: data.OrderNo,
    })
    return err
}

func (h *InventoryHandler) handleOrderPaid(event *types.OrderEvent) error {
    data, ok := event.Data.(*types.OrderPaidData)
    if !ok {
        return fmt.Errorf("invalid event data type")
    }

    _, err := h.inventoryRpc.DeductStock(context.Background(), &inventoryclient.DeductStockRequest{
        OrderNo: data.OrderNo,
    })
    return err
}