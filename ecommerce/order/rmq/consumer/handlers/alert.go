package handlers

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/streadway/amqp"
)

type AlertHandler struct {
    logger         *zerolog.Logger
    ordersModel    model.OrdersModel
    paymentsModel  model.OrderPaymentsModel
    refundsModel   model.OrderRefundsModel
}

func NewAlertHandler(
    ordersModel model.OrdersModel,
    paymentsModel model.OrderPaymentsModel,
    refundsModel model.OrderRefundsModel,
) *AlertHandler {
    return &AlertHandler{
        logger:         zerolog.GetLogger(),
        ordersModel:    ordersModel,
        paymentsModel:  paymentsModel,
        refundsModel:   refundsModel,
    }
}

func (h *AlertHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.OrderAlertEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "order_no":    event.OrderNo,
        "alert_type":  event.AlertType,
        "alert_level": event.AlertLevel,
        "message":     event.Message,
    }
    h.logger.Info(ctx, "Processing order alert event", fields)

    switch event.AlertType {
    case "payment_timeout":
        return h.handlePaymentTimeout(ctx, event)
    case "shipping_delay":
        return h.handleShippingDelay(ctx, event)
    case "refund_request":
        return h.handleRefundRequest(ctx, event)
    }

    return nil
}

func (h *AlertHandler) handlePaymentTimeout(ctx context.Context, event types.OrderAlertEvent) error {
    // Cancel order due to payment timeout
    order, err := h.ordersModel.FindByOrderNo(ctx, event.OrderNo)
    if err != nil {
        return err
    }

    order.Status = 5 // Canceled
    order.Remark = sql.NullString{String: "Payment timeout", Valid: true}
    
    return h.ordersModel.Update(ctx, order)
}

func (h *AlertHandler) handleShippingDelay(ctx context.Context, event types.OrderAlertEvent) error {
    // Log shipping delay alert
    h.logger.Warn(ctx, "Shipping delay detected", map[string]interface{}{
        "order_no": event.OrderNo,
        "message": event.Message,
    })
    return nil
}

func (h *AlertHandler) handleRefundRequest(ctx context.Context, event types.OrderAlertEvent) error {
    // Update order status to refunding
	order, err := h.ordersModel.FindByOrderNo(ctx, event.OrderNo)
	if err != nil {
		return err 
	}
	
    return h.ordersModel.UpdateStatus(ctx, order.Id, 6) // 6: Refunding
}