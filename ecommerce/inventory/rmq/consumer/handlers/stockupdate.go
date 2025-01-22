package handlers

import (
	"context"
	"fmt"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/middleware"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventoryclient"
)

type StockUpdateHandler struct {
	BaseHandler
	inventoryRpc inventoryclient.Inventory
}

func NewStockUpdateHandler(
	logger middleware.Logger,
	inventoryRpc inventoryclient.Inventory,
) *StockUpdateHandler {
	return &StockUpdateHandler{
		BaseHandler:  NewBaseHandler(logger),
		inventoryRpc: inventoryRpc,
	}
}

func (h *StockUpdateHandler) Handle(event *types.InventoryEvent) error {
	h.LogEvent(event, "handling stock update event")

	data, ok := event.Data.(*types.StockUpdateData)
	if !ok {
		return fmt.Errorf("invalid event data type")
	}

	req := &inventoryclient.UpdateStockRequest{
		SkuId:       int64(data.SkuID),
		WarehouseId: int64(data.WarehouseID),
		Quantity:    data.Quantity,
		Remark:      data.Remark,
	}

	_, err := h.inventoryRpc.UpdateStock(context.Background(), req)
	if err != nil {
		h.LogError(event, err)
		return types.NewRetryableError(err)
	}

	h.LogEvent(event, "stock update event handled successfully")
	return nil
}
