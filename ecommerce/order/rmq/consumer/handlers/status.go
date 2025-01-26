package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/streadway/amqp"
)

type StatusHandler struct {
    logger         *zerolog.Logger
    ordersModel    model.OrdersModel
    paymentsModel  model.OrderPaymentsModel
    shippingModel  model.OrderShippingModel
    refundsModel   model.OrderRefundsModel
}

func NewStatusHandler(
    ordersModel model.OrdersModel,
    paymentsModel model.OrderPaymentsModel,
    shippingModel model.OrderShippingModel,
    refundsModel model.OrderRefundsModel,
) *StatusHandler {
    return &StatusHandler{
        logger:         zerolog.GetLogger(),
        ordersModel:    ordersModel,
        paymentsModel:  paymentsModel,
        shippingModel:  shippingModel,
        refundsModel:   refundsModel,
    }
}

func (h *StatusHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.OrderStatusChangedEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "order_no":   event.OrderNo,
        "user_id":    event.UserID,
        "old_status": event.OldStatus,
        "new_status": event.NewStatus,
        "event_type": event.EventType,
    }
    h.logger.Info(ctx, "Processing order status change event", fields)

    switch event.EventType {
    case types.OrderStatusPaid:
        return h.handleOrderPaid(ctx, event)
    case types.OrderStatusShipped:
        return h.handleOrderShipped(ctx, event)
    case types.OrderStatusReceived:
        return h.handleOrderReceived(ctx, event)
    case types.OrderStatusCanceled:
        return h.handleOrderCanceled(ctx, event)
    case types.OrderStatusRefunding:
        return h.handleOrderRefunding(ctx, event)
    default:
        return h.updateOrderStatus(ctx, event)
    }
}

func (h *StatusHandler) handleOrderPaid(ctx context.Context, event types.OrderStatusChangedEvent) error {
    // Update payment status
    if err := h.paymentsModel.UpdateStatus(ctx, event.PaymentNo, 1, time.Now()); err != nil {
        return err
    }
    
    // Update order status
    return h.ordersModel.UpdateStatus(ctx, uint64(event.UserID), 2) // 2: Waiting for shipment
}

func (h *StatusHandler) handleOrderShipped(ctx context.Context, event types.OrderStatusChangedEvent) error {
    // Update shipping status
	order, err := h.ordersModel.FindByOrderNo(ctx, event.OrderNo)
	if err != nil {
		return err 
	}
    shipping, err := h.shippingModel.FindByOrderId(ctx, order.Id)
    if err != nil {
        return err
    }
    
    if err := h.shippingModel.UpdateStatus(ctx, shipping.Id, 1); err != nil {
        return err
    }
    
    // Update order status
    return h.ordersModel.UpdateStatus(ctx, uint64(event.UserID), 3) // 3: Waiting for receipt
}

func (h *StatusHandler) handleOrderReceived(ctx context.Context, event types.OrderStatusChangedEvent) error {
    // Update shipping status
	order, err := h.ordersModel.FindByOrderNo(ctx, event.OrderNo)
	if err != nil {
		return err 
	}

    shipping, err := h.shippingModel.FindByOrderId(ctx, order.Id)
    if err != nil {
        return err
    }
    
    if err := h.shippingModel.UpdateStatus(ctx, shipping.Id, 2); err != nil {
        return err
    }
    
    // Update order status
    return h.ordersModel.UpdateStatus(ctx, uint64(event.UserID), 4) // 4: Completed
}

func (h *StatusHandler) handleOrderCanceled(ctx context.Context, event types.OrderStatusChangedEvent) error {
    // Update order status and add cancel reason
    order, err := h.ordersModel.FindByOrderNo(ctx, event.OrderNo)
    if err != nil {
        return err
    }
    
    order.Status = 5 // 5: Canceled
    order.Remark = sql.NullString{String: event.Reason, Valid: true}
    
    return h.ordersModel.Update(ctx, order)
}

func (h *StatusHandler) handleOrderRefunding(ctx context.Context, event types.OrderStatusChangedEvent) error {
    // Update order status
    return h.ordersModel.UpdateStatus(ctx, uint64(event.UserID), 6) // 6: Refunding
}

func (h *StatusHandler) updateOrderStatus(ctx context.Context, event types.OrderStatusChangedEvent) error {
    return h.ordersModel.UpdateStatus(ctx, uint64(event.UserID), int64(event.NewStatus))
}