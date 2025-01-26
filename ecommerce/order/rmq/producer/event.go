package producer

import (
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
)

func NewOrderStatusChangedEvent(
    orderNo string,
    userId int64,
    oldStatus int32,
    newStatus int32,
    eventType types.OrderStatusEventType,
    paymentNo string,
    shippingNo string,
    refundNo string,
    reason string,
) *types.OrderStatusChangedEvent {
    return &types.OrderStatusChangedEvent{
        OrderEvent: types.OrderEvent{
            Type:      types.OrderEventType(eventType),
            OrderNo:   orderNo,
            UserID:    userId,
            Timestamp: time.Now(),
        },
        OldStatus:  oldStatus,
        NewStatus:  newStatus,
        EventType:  eventType,
        PaymentNo:  paymentNo,
        ShippingNo: shippingNo,
        RefundNo:   refundNo,
        Reason:     reason,
    }
}

func NewOrderAlertEvent(
    orderNo string,
    userId int64,
    alertType string,
    alertLevel string,
    message string,
) *types.OrderAlertEvent {
    return &types.OrderAlertEvent{
        OrderEvent: types.OrderEvent{
            Type:      types.OrderEventType("order.alert"),
            OrderNo:   orderNo,
            UserID:    userId,
            Timestamp: time.Now(),
        },
        AlertType:  alertType,
        AlertLevel: alertLevel,
        Message:    message,
    }
}