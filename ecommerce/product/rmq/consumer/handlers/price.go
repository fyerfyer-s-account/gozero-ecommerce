package handlers

import (
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rmq/types"
)

type PriceHandler struct {
	// userRpc user.UserClient
	// orderRpc order.OrderClient
}

func NewPriceHandler() *PriceHandler {
	return &PriceHandler{}
}

func (h *PriceHandler) HandlePriceUpdate(event *types.ProductEvent) error {
	_, ok := event.Data.(*types.PriceData)
	if !ok {
		return zeroerr.ErrInvalidEventData
	}

	// TODO: Notify users about price drops
	// TODO: Update order service about price changes
	return nil
}
