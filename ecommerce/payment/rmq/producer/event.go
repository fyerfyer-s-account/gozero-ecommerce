package producer

import (
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
)

func NewPaymentCreatedEvent(
    orderNo string,
    paymentNo string,
    amount float64,
    paymentMethod int32,
    payURL string,
) *types.PaymentCreatedEvent {
    return &types.PaymentCreatedEvent{
        PaymentEvent: types.PaymentEvent{
            Type:      types.PaymentCreated,
            OrderNo:   orderNo,
            PaymentNo: paymentNo,
            Timestamp: time.Now(),
        },
        Amount:        amount,
        PaymentMethod: paymentMethod,
        PayURL:       payURL,
    }
}

func NewPaymentSuccessEvent(
    orderNo string,
    paymentNo string,
    amount float64,
    paymentMethod int32,
    paidTime time.Time,
) *types.PaymentSuccessEvent {
    return &types.PaymentSuccessEvent{
        PaymentEvent: types.PaymentEvent{
            Type:      types.PaymentSuccess,
            OrderNo:   orderNo,
            PaymentNo: paymentNo,
            Timestamp: time.Now(),
        },
        Amount:        amount,
        PaymentMethod: paymentMethod,
        PaidTime:     paidTime,
    }
}

func NewPaymentFailedEvent(
    orderNo string,
    paymentNo string,
    amount float64,
    reason string,
    errorCode string,
) *types.PaymentFailedEvent {
    return &types.PaymentFailedEvent{
        PaymentEvent: types.PaymentEvent{
            Type:      types.PaymentFailed,
            OrderNo:   orderNo,
            PaymentNo: paymentNo,
            Timestamp: time.Now(),
        },
        Amount:    amount,
        Reason:    reason,
        ErrorCode: errorCode,
    }
}

func NewPaymentRefundEvent(
    orderNo string,
    paymentNo string,
    refundNo string,
    refundAmount float64,
    reason string,
    refundTime time.Time,
) *types.PaymentRefundEvent {
    return &types.PaymentRefundEvent{
        PaymentEvent: types.PaymentEvent{
            Type:      types.PaymentRefund,
            OrderNo:   orderNo,
            PaymentNo: paymentNo,
            Timestamp: time.Now(),
        },
        RefundNo:     refundNo,
        RefundAmount: refundAmount,
        Reason:       reason,
        RefundTime:   refundTime,
    }
}

func NewPaymentVerificationEvent(
    orderNo string,
    paymentNo string,
    verified bool,
    message string,
) *types.PaymentVerificationEvent {
    return &types.PaymentVerificationEvent{
        PaymentEvent: types.PaymentEvent{
            Type:      types.PaymentVerified,
            OrderNo:   orderNo,
            PaymentNo: paymentNo,
            Timestamp: time.Now(),
        },
        Verified: verified,
        Message:  message,
    }
}