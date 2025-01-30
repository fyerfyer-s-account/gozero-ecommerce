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

    fields := map[string]interface{}{
        "order_no":    event.OrderNo,
        "payment_no":  event.PaymentNo,
        "amount":      event.Amount,
        "reason":      event.Reason,
        "error_code":  event.ErrorCode,
    }
    h.logger.Info(ctx, "Processing payment failed event", fields)

    // Get order
    order, err := h.ordersModel.FindByOrderNo(ctx, event.OrderNo)
    if err != nil {
        h.logger.Error(ctx, "Failed to find order", err, fields)
        return err
    }

    // Update order status to cancelled (5)
    if err := h.ordersModel.UpdateStatus(ctx, order.Id, 5); err != nil {
        h.logger.Error(ctx, "Failed to update order status", err, fields)
        return err
    }

    // Get payment record
    payment, err := h.paymentsModel.FindOneByPaymentNo(ctx, event.PaymentNo)
    if err != nil {
        h.logger.Error(ctx, "Failed to find payment", err, fields)
        return err
    }

    // Update payment status to failed
    if err := h.paymentsModel.UpdateStatus(ctx, event.PaymentNo, 0, payment.PayTime.Time); err != nil {
        h.logger.Error(ctx, "Failed to update payment status", err, fields)
        return err
    }

    h.logger.Info(ctx, "Successfully processed payment failed event", fields)
    return nil
}