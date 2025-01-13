package handlers

import (
    "context"
    "fmt"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventoryclient"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/types"
)

type InventoryHandler struct {
    inventoryRpc inventoryclient.Inventory
}

func NewInventoryHandler(inventoryRpc inventoryclient.Inventory) *InventoryHandler {
    return &InventoryHandler{
        inventoryRpc: inventoryRpc,
    }
}

func (h *InventoryHandler) Handle(event *types.OrderEvent) error {
    switch event.Type {
    case types.EventTypeOrderCreated:
        return h.handleOrderCreated(event)
    case types.EventTypeOrderCancelled:
        return h.handleOrderCancelled(event)
    case types.EventTypeOrderPaid:
        return h.handleOrderPaid(event)
    default:
        return nil
    }
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