package handlers

import (
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rmq/types"
)

type AddressHandler struct {
	// orderRpc order.OrderClient
}

func NewAddressHandler() *AddressHandler {
	return &AddressHandler{}
}

func (h *AddressHandler) HandleAddressUpdate(event *types.UserEvent) error {
	_, ok := event.Data.(*types.AddressData)
	if !ok {
		return zeroerr.ErrInvalidEventData
	}

	// TODO: Notify order service about address changes
	return nil
}
