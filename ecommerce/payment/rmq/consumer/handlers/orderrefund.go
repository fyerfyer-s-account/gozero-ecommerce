package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/util"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/streadway/amqp"
)

type OrderRefundHandler struct {
    logger          *zerolog.Logger
    paymentOrders   model.PaymentOrdersModel
    refundOrders    model.RefundOrdersModel
    paymentLogs     model.PaymentLogsModel
}

func NewOrderRefundHandler(
    paymentOrders model.PaymentOrdersModel,
    refundOrders model.RefundOrdersModel,
    paymentLogs model.PaymentLogsModel,
) *OrderRefundHandler {
    return &OrderRefundHandler{
        logger:          zerolog.GetLogger(),
        paymentOrders:   paymentOrders,
        refundOrders:    refundOrders,
        paymentLogs:     paymentLogs,
    }
}

func (h *OrderRefundHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.OrderRefundedEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "order_no":    event.OrderNo,
        "user_id":     event.UserID,
        "refund_type": event.RefundType,
        "amount":      event.Amount,
        "timestamp":   event.Timestamp,
    }
    h.logger.Info(ctx, "Processing order refund request", fields)

    // Get payment order
    payment, err := h.paymentOrders.FindByOrderNo(ctx, event.OrderNo)
    if err != nil {
        h.logger.Error(ctx, "Failed to find payment order", err, fields)
        return err
    }
    if len(payment) == 0 {
        return fmt.Errorf("payment order not found for order: %s", event.OrderNo)
    }

    // Generate refund number
    refundNo := util.GenerateNo("REF")

    // Create refund order
    refund := &model.RefundOrders{
        RefundNo:    refundNo,
        PaymentNo:   payment[0].PaymentNo,
        OrderNo:     event.OrderNo,
        UserId:      uint64(event.UserID),
        Amount:      event.Amount,
        Reason:      event.Reason,
        Status:      1, // 1: Pending
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    // Insert refund order
    _, err = h.refundOrders.Insert(ctx, refund)
    if err != nil {
        h.logger.Error(ctx, "Failed to create refund order", err, fields)
        return err
    }

    data, err := json.Marshal(event)
    if err != nil {
        h.logger.Error(ctx, "Failed to marshal event data", err, fields)
    }

    // Log refund request
    log := &model.PaymentLogs{
        PaymentNo:   payment[0].PaymentNo,
        Type:        2, // 2: Refund
        Channel:     payment[0].Channel,
        RequestData: sql.NullString{
            String: string(data),
            Valid:  string(data) != "",
        },
        CreatedAt:   time.Now(),
    }

    if _, err := h.paymentLogs.Insert(ctx, log); err != nil {
        h.logger.Error(ctx, "Failed to create payment log", err, fields)
        // Don't return error here, continue processing
    }

    // Update payment order status
    if err := h.paymentOrders.UpdateStatus(ctx, payment[0].Id, 4); err != nil { // 4: Refunded
        h.logger.Error(ctx, "Failed to update payment status", err, fields)
        return err
    }

    err = msg.Ack(false)
    if err != nil {
        h.logger.Error(ctx, "Failed to acknowledge message", err, fields)
        return err
    }

    h.logger.Info(ctx, "Refund order created successfully", map[string]interface{}{
        "refund_no":  refundNo,
        "payment_no": payment[0].PaymentNo,
        "order_no":   event.OrderNo,
        "amount":     event.Amount,
    })

    return nil
}