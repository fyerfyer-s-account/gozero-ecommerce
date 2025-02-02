package handlers

import (
	"context"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/streadway/amqp"
)

type PaymentFailedHandler struct {
    logger         *zerolog.Logger
    ordersModel    model.OrdersModel
    paymentsModel  model.OrderPaymentsModel
}

func NewPaymentFailedHandler(
    ordersModel model.OrdersModel,
    paymentsModel model.OrderPaymentsModel,
) *PaymentFailedHandler {
    return &PaymentFailedHandler{
        logger:         zerolog.GetLogger(),
        ordersModel:    ordersModel,
        paymentsModel:  paymentsModel,
    }
}

func (h *PaymentFailedHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.OrderPaymentFailedEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }

    h.logger.Info(ctx, "Received payment failed event", map[string]interface{}{
        "order_no":    event.OrderNo,
        "payment_no":  event.PaymentNo,
        "amount":      event.Amount,
        "reason":      event.Reason,
        "error_code":  event.ErrorCode,
    })

    // 1. Update order status
    order, err := h.ordersModel.FindByOrderNo(ctx, event.OrderNo)
    if err != nil {
        h.logger.Error(ctx, "Failed to find order", err, nil)
        return err
    }

    if err := h.ordersModel.UpdateStatus(ctx, order.Id, 5); err != nil {
        h.logger.Error(ctx, "Failed to update order status", err, nil)
        return err
    }

    // 2. Update payment status
    payment, err := h.paymentsModel.FindOneByPaymentNo(ctx, event.PaymentNo)
    if err != nil {
        h.logger.Error(ctx, "Failed to find payment", err, nil)
        return err
    }

    if err := h.paymentsModel.UpdateStatus(ctx, payment.PaymentNo, 0, payment.CreatedAt); err != nil {
        h.logger.Error(ctx, "Failed to update payment status", err, nil)
        return err
    }

    // Acknowledge message
    if err := msg.Ack(false); err != nil {
        h.logger.Error(ctx, "Failed to acknowledge message", err, nil)
        return err
    }

    h.logger.Info(ctx, "Successfully processed payment failed event", nil)
    return nil
}