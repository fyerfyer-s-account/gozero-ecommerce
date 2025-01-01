package handlers

import (
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rmq/types"
)

type WalletHandler struct {
	// paymentRpc payment.PaymentClient
}

func NewWalletHandler() *WalletHandler {
	return &WalletHandler{}
}

func (h *WalletHandler) HandleTransaction(event *types.UserEvent) error {
	_, ok := event.Data.(*types.WalletData)
	if !ok {
		return zeroerr.ErrInvalidEventData
	}

	// TODO: Notify payment service about balance changes
	// TODO: Send transaction notifications
	return nil
}
