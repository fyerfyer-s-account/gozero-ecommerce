package handlers

import (
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rmq/types"
)

type ProfileHandler struct {
	// RPC clients will be added later
	// productRpc    product.ProductClient
	// orderRpc     order.OrderClient
}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{}
}

func (h *ProfileHandler) HandleProfileUpdate(event *types.UserEvent) error {
	_, ok := event.Data.(*types.ProfileData)
	if !ok {
		return zeroerr.ErrInvalidEventData
	}

	// TODO: Notify other services about user preference changes
	// TODO: Send email/SMS notifications
	return nil
}
