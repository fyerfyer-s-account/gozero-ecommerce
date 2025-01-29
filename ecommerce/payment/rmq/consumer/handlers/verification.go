package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/streadway/amqp"
)

type PaymentVerificationHandler struct {
    logger          *zerolog.Logger
    paymentOrders   model.PaymentOrdersModel
    paymentChannels model.PaymentChannelsModel
    paymentLogs     model.PaymentLogsModel
}

func NewPaymentVerificationHandler(
    paymentOrders model.PaymentOrdersModel,
    paymentChannels model.PaymentChannelsModel,
    paymentLogs model.PaymentLogsModel,
) *PaymentVerificationHandler {
    return &PaymentVerificationHandler{
        logger:          zerolog.GetLogger(),
        paymentOrders:   paymentOrders,
        paymentChannels: paymentChannels,
        paymentLogs:     paymentLogs,
    }
}

func (h *PaymentVerificationHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.PaymentEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "payment_no": event.PaymentNo,
        "order_no":   event.OrderNo,
        "timestamp":  event.Timestamp,
    }
    h.logger.Info(ctx, "Processing payment verification request", fields)

    // Get payment order
    payment, err := h.paymentOrders.FindOneByPaymentNo(ctx, event.PaymentNo)
    if err != nil {
        h.logger.Error(ctx, "Failed to find payment order", err, fields)
        return err
    }

    // Get payment channel
    channel, err := h.paymentChannels.FindOneByChannel(ctx, payment.Channel)
    if err != nil {
        h.logger.Error(ctx, "Failed to find payment channel", err, fields)
        return err
    }

    // Verify payment status with channel
    verified := h.verifyPaymentStatus(ctx, payment, channel)

    // Create verification log
	reqData, err := json.Marshal(event)
	if err != nil {
		h.logger.Error(ctx, "Failed to marshal event data", err, fields)
		return err 
	}

	resData, err := json.Marshal(event)
	if err != nil {
		h.logger.Error(ctx, "Failed to marshal validate data", err, fields)
		return err 
	}

    verifyLog := &model.PaymentLogs{
        PaymentNo:    payment.PaymentNo,
        Type:        3, // 3: Verification
        Channel:     payment.Channel,
        RequestData: sql.NullString{
            String: string(reqData),
            Valid:  string(reqData) != "",
        },
        ResponseData: sql.NullString{
            String: string(resData),
            Valid:  string(resData) != "",
        },
        CreatedAt:   time.Now(),
    }

    if _, err := h.paymentLogs.Insert(ctx, verifyLog); err != nil {
        h.logger.Error(ctx, "Failed to create verification log", err, fields)
        return err
    }

    // Update payment status if verified
    if verified.Verified {
        updates := map[string]interface{}{
            "status":     3, // 3: Paid
            "pay_time":   sql.NullTime{Time: time.Now(), Valid: true},
            "channel_data": sql.NullString{String: verified.Message, Valid: true},
        }
        
        if err := h.paymentOrders.UpdatePartial(ctx, payment.Id, updates); err != nil {
            h.logger.Error(ctx, "Failed to update payment status", err, fields)
            return err
        }
    }

    // Publish verification result event
    verificationEvent := types.PaymentVerificationEvent{
        PaymentEvent: event,
        Verified:    verified.Verified,
        Message:     verified.Message,
    }

    _, err = json.Marshal(verificationEvent)
    if err != nil {
        return err
    }

    err = msg.Ack(false)
    if err != nil {
        h.logger.Error(ctx, "Failed to acknowledge message", err, fields)
        return err
    }

    h.logger.Info(ctx, "Payment verification completed", map[string]interface{}{
        "payment_no": payment.PaymentNo,
        "verified":   verified.Verified,
        "message":    verified.Message,
    })

    return nil
}

type verificationResult struct {
    Verified bool   `json:"verified"`
    Message  string `json:"message"`
}

func (h *PaymentVerificationHandler) verifyPaymentStatus(
    ctx context.Context, 
    payment *model.PaymentOrders, 
    channel *model.PaymentChannels,
) *verificationResult {
    // Mock verification logic - in real world this would call the payment provider's API
    if payment.Status == 2 { // Payment in progress
        return &verificationResult{
            Verified: true,
            Message: "Payment verified successfully",
        }
    }

    return &verificationResult{
        Verified: false,
        Message: "Payment not verified",
    }
}