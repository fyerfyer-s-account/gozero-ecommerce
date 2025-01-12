package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/paymentclient"
)

type PaymentHandler struct {
	paymentRpc         paymentclient.Payment
	orderModel         model.OrdersModel
	orderRefundsModel  model.OrderRefundsModel
	orderPaymentsModel model.OrderPaymentsModel
}

func NewPaymentHandler(paymentRpc paymentclient.Payment, orderModel model.OrdersModel,
	orderRefundsModel model.OrderRefundsModel, orderPaymentsModel model.OrderPaymentsModel) *PaymentHandler {
	return &PaymentHandler{
		paymentRpc:         paymentRpc,
		orderModel:         orderModel,
		orderRefundsModel:  orderRefundsModel,
		orderPaymentsModel: orderPaymentsModel,
	}
}

func (h *PaymentHandler) Handle(event *types.OrderEvent) error {
	switch event.Type {
	case types.EventTypeOrderPaid:
		return h.handleOrderPaid(event)
	case types.EventTypeOrderCancelled:
		return h.handleOrderCancelled(event)
	default:
		return nil
	}
}

func (h *PaymentHandler) handleOrderPaid(event *types.OrderEvent) error {
	data, ok := event.Data.(*types.OrderPaidData)
	if !ok {
		return fmt.Errorf("invalid event data type")
	}

	// 1. Update payment status
	payment, err := h.orderPaymentsModel.FindOneByPaymentNo(context.Background(), data.PaymentNo)
	if err != nil {
		return err
	}

	err = h.orderPaymentsModel.UpdateStatus(context.Background(), data.PaymentNo,
		1, data.PayTime)
	if err != nil {
		return err
	}

	// 2. Update order status
	return h.orderModel.UpdateStatus(context.Background(), payment.OrderId, 4)
}

func (h *PaymentHandler) handleOrderCancelled(event *types.OrderEvent) error {
    data, ok := event.Data.(*types.OrderCancelledData)
    if !ok {
        return fmt.Errorf("invalid event data type")
    }

    refund := &model.OrderRefunds{
        OrderId:     uint64(data.OrderId),
        RefundNo:    fmt.Sprintf("R%d", time.Now().UnixNano()),
        Amount:      data.Amount,
        Reason:      data.Reason,
        Status:      0,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    _, err := h.orderRefundsModel.Insert(context.Background(), refund)
    return err
}
