package handlers

import (
    "context"
    "encoding/json"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type PaymentSuccessHandler struct {
    logger         *zerolog.Logger
    ordersModel    model.OrdersModel
    paymentsModel  model.OrderPaymentsModel
}

func NewPaymentSuccessHandler(
    ordersModel model.OrdersModel,
    paymentsModel model.OrderPaymentsModel,
) *PaymentSuccessHandler {
    return &PaymentSuccessHandler{
        logger:         zerolog.GetLogger(),
        ordersModel:    ordersModel,
        paymentsModel:  paymentsModel,
    }
}

func (h *PaymentSuccessHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.OrderPaymentSuccessEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "order_no":       event.OrderNo,
        "payment_no":     event.PaymentNo,
        "amount":        event.Amount,
        "payment_method": event.PaymentMethod,
        "paid_time":     event.PaidTime,
    }
    h.logger.Info(ctx, "Processing payment success event", fields)

    // Get order
    order, err := h.ordersModel.FindByOrderNo(ctx, event.OrderNo)
    if err != nil {
        h.logger.Error(ctx, "Failed to find order", err, fields)
        return err
    }

    // Update order status to paid (2: Pending shipment)
    if err := h.ordersModel.UpdateStatus(ctx, order.Id, 2); err != nil {
        h.logger.Error(ctx, "Failed to update order status", err, fields)
        return err
    }

    // Update payment record
    _, err = h.paymentsModel.FindOneByPaymentNo(ctx, event.PaymentNo)
    if err != nil {
        h.logger.Error(ctx, "Failed to find payment", err, fields)
        return err
    }

    // Update payment status to paid (1)
    if err := h.paymentsModel.UpdateStatus(ctx, event.PaymentNo, 1, event.PaidTime); err != nil {
        h.logger.Error(ctx, "Failed to update payment status", err, fields)
        return err
    }

    h.logger.Info(ctx, "Successfully processed payment success event", fields)
    return nil
}