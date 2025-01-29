package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/util"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/streadway/amqp"
)

type OrderPaymentHandler struct {
	logger          *zerolog.Logger
	paymentOrders   model.PaymentOrdersModel
	paymentChannels model.PaymentChannelsModel
	paymentLogs     model.PaymentLogsModel
}

func NewOrderPaymentHandler(
	paymentOrders model.PaymentOrdersModel,
	paymentChannels model.PaymentChannelsModel,
	paymentLogs model.PaymentLogsModel,
) *OrderPaymentHandler {
	return &OrderPaymentHandler{
		logger:          zerolog.GetLogger(),
		paymentOrders:   paymentOrders,
		paymentChannels: paymentChannels,
		paymentLogs:     paymentLogs,
	}
}

func (h *OrderPaymentHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
	var event types.OrderEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return err
	}

	fields := map[string]interface{}{
		"order_no":  event.OrderNo,
		"user_id":   event.UserID,
		"timestamp": event.Timestamp,
	}
	h.logger.Info(ctx, "Processing order payment request", fields)

	// Cast to OrderPaidEvent to get payment details
	var paidEvent types.OrderPaidEvent
	if err := json.Unmarshal(msg.Body, &paidEvent); err != nil {
		return err
	}

	// Generate payment number
	paymentNo := util.GenerateNo("PAY")

	// Create payment order
	payment := &model.PaymentOrders{
		PaymentNo: paymentNo,
		OrderNo:   event.OrderNo,
		UserId:    uint64(event.UserID),
		Amount:    paidEvent.PayAmount,
		Channel:   int64(paidEvent.PaymentMethod),
		Status:    1, // 1: Pending
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Validate payment channel
	channel, err := h.paymentChannels.FindOneByChannelAndStatus(ctx, payment.Channel, 1)
	if err != nil {
		h.logger.Error(ctx, "Payment channel not found or disabled", err, fields)
		return err
	}

	// Create payment order record
	_, err = h.paymentOrders.Insert(ctx, payment)
	if err != nil {
		h.logger.Error(ctx, "Failed to create payment order", err, fields)
		return err
	}

	// Log payment request
	log := &model.PaymentLogs{
		PaymentNo: paymentNo,
		Type:      1, // 1: Payment
		Channel:   payment.Channel,
		RequestData: sql.NullString{
			String: channel.Config,
			Valid:  channel.Config != "",
		}, // Use channel config as request data
		CreatedAt: time.Now(),
	}

	_, err = h.paymentLogs.Insert(ctx, log)
	if err != nil {
		h.logger.Error(ctx, "Failed to create payment log", err, fields)
		return err
	}

	// Publish payment created event
	createdEvent := types.PaymentEvent{
		Type:      types.PaymentCreated,
		OrderNo:   event.OrderNo,
		PaymentNo: paymentNo,
		Timestamp: time.Now(),
	}

	_, err = json.Marshal(createdEvent)
	if err != nil {
		return err
	}

	err = msg.Ack(false)
	if err != nil {
		h.logger.Error(ctx, "Failed to acknowledge message", err, fields)
		return err
	}

	h.logger.Info(ctx, "Payment order created successfully", map[string]interface{}{
		"payment_no": paymentNo,
		"order_no":   event.OrderNo,
		"amount":     payment.Amount,
		"channel":    payment.Channel,
	})

	return nil
}
