package handlers

import (
    "context"
    "encoding/json"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type PaymentRefundHandler struct {
    logger         *zerolog.Logger
    ordersModel    model.OrdersModel
    paymentsModel  model.OrderPaymentsModel
    refundsModel   model.OrderRefundsModel
}

func NewPaymentRefundHandler(
    ordersModel model.OrdersModel,
    paymentsModel model.OrderPaymentsModel,
    refundsModel model.OrderRefundsModel,
) *PaymentRefundHandler {
    return &PaymentRefundHandler{
        logger:         zerolog.GetLogger(),
        ordersModel:    ordersModel,
        paymentsModel:  paymentsModel,
        refundsModel:   refundsModel,
    }
}

func (h *PaymentRefundHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.OrderPaymentRefundedEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "order_no":      event.OrderNo,
        "payment_no":    event.PaymentNo,
        "refund_no":     event.RefundNo,
        "refund_amount": event.RefundAmount,
        "reason":        event.Reason,
        "refund_time":   event.RefundTime,
    }
    h.logger.Info(ctx, "Processing payment refund event", fields)

    // Get order
    order, err := h.ordersModel.FindByOrderNo(ctx, event.OrderNo)
    if err != nil {
        h.logger.Error(ctx, "Failed to find order", err, fields)
        return err
    }

    // Update order status to refunded
    if err := h.ordersModel.UpdateStatus(ctx, order.Id, 6); err != nil { // 6: Refunded
        h.logger.Error(ctx, "Failed to update order status", err, fields)
        return err
    }

    // Update payment status
    _, err = h.paymentsModel.FindOneByPaymentNo(ctx, event.PaymentNo)
    if err != nil {
        h.logger.Error(ctx, "Failed to find payment", err, fields)
        return err
    }

    if err := h.paymentsModel.UpdateStatus(ctx, event.PaymentNo, 2, event.RefundTime); err != nil { // 2: Refunded
        h.logger.Error(ctx, "Failed to update payment status", err, fields)
        return err
    }

    // Update refund status if exists
    _, err = h.refundsModel.FindOneByRefundNo(ctx, event.RefundNo)
    if err == nil {
        if err := h.refundsModel.UpdateStatus(ctx, event.RefundNo, 3, "Refund completed"); err != nil { // 3: Refunded
            h.logger.Error(ctx, "Failed to update refund status", err, fields)
            return err
        }
    }

    h.logger.Info(ctx, "Successfully processed payment refund event", fields)
    return nil
}