package handlers

import (
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rmq/types"
)

type InventoryHandler struct {
	// orderRpc order.OrderClient
	// userRpc user.UserClient
}

func NewInventoryHandler() *InventoryHandler {
	return &InventoryHandler{}
}

func (h *InventoryHandler) HandleStockUpdate(event *types.ProductEvent) error {
	_, ok := event.Data.(*types.StockData)
	if !ok {
		return zeroerr.ErrInvalidEventData
	}

	// TODO: Notify order service about stock changes
	// TODO: Send back-in-stock notifications to users
	return nil
}
