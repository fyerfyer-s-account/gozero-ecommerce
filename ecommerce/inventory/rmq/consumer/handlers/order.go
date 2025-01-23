package handlers

import (
	"context"
	"fmt"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/middleware"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventoryclient"
)

type OrderEventHandler struct {
	BaseHandler
	inventoryRpc inventoryclient.Inventory
}

func NewOrderEventHandler(logger middleware.Logger, inventoryRpc inventoryclient.Inventory) *OrderEventHandler {
	return &OrderEventHandler{
		BaseHandler:  NewBaseHandler(logger),
		inventoryRpc: inventoryRpc,
	}
}

func (h *OrderEventHandler) Handle(event *types.InventoryEvent) error {
	h.LogEvent(event, "handling order event")

	var err error
	switch event.Type {
	case types.EventTypeOrderCreated:
		err = h.handleOrderCreated(event)
	case types.EventTypeOrderCancelled:
		err = h.handleOrderCancelled(event)
	case types.EventTypeOrderPaid:
		err = h.handleOrderPaid(event)
	case types.EventTypeOrderRefunded:
		err = h.handleOrderRefunded(event)
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}

	if err != nil {
		h.LogError(event, err)
		return types.NewRetryableError(err)
	}

	h.LogEvent(event, "order event handled successfully")
	return nil
}

func (h *OrderEventHandler) handleOrderCreated(event *types.InventoryEvent) error {
	data, ok := event.Data.(*types.OrderCreatedData)
	if !ok {
		return fmt.Errorf("invalid event data type")
	}

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
		Items:   items,
	})
	return err
}

func (h *OrderEventHandler) handleOrderCancelled(event *types.InventoryEvent) error {
	data, ok := event.Data.(*types.OrderCancelledData)
	if !ok {
		return fmt.Errorf("invalid event data type")
	}

	_, err := h.inventoryRpc.UnlockStock(context.Background(), &inventoryclient.UnlockStockRequest{
		OrderNo: data.OrderNo,
	})
	return err
}

func (h *OrderEventHandler) handleOrderPaid(event *types.InventoryEvent) error {
	data, ok := event.Data.(*types.OrderPaidData)
	if !ok {
		return fmt.Errorf("invalid event data type")
	}

	_, err := h.inventoryRpc.DeductStock(context.Background(), &inventoryclient.DeductStockRequest{
		OrderNo: data.OrderNo,
	})
	return err
}

func (h *OrderEventHandler) handleOrderRefunded(event *types.InventoryEvent) error {
	data, ok := event.Data.(*types.OrderRefundedData)
	if !ok {
		return fmt.Errorf("invalid event data type")
	}

	// Return stock to inventory
	for _, item := range data.Items {
		_, err := h.inventoryRpc.UpdateStock(context.Background(), &inventoryclient.UpdateStockRequest{
			SkuId:       int64(item.SkuID),
			WarehouseId: int64(item.WarehouseID),
			Quantity:    item.Quantity,
			Remark:      fmt.Sprintf("Refund order: %s", data.OrderNo),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
