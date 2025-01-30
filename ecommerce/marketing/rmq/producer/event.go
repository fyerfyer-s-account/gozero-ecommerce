package producer

import (
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
)

func NewCouponEvent(
    userId int64,
    couponId int64,
    couponCode string,
    amount float64,
    orderNo string,
    eventType types.MarketingEventType,
) *types.CouponEvent {
    return &types.CouponEvent{
        MarketingEvent: types.MarketingEvent{
            Type:      eventType,
            UserID:    userId,
            Timestamp: time.Now(),
        },
        CouponID:   couponId,
        CouponCode: couponCode,
        Amount:     amount,
        OrderNo:    orderNo,
    }
}

func NewPromotionEvent(
    userId int64,
    promotionId int64,
    promotionType int32,
    discount float64,
    orderNo string,
    eventType types.MarketingEventType,
) *types.PromotionEvent {
    return &types.PromotionEvent{
        MarketingEvent: types.MarketingEvent{
            Type:      eventType,
            UserID:    userId,
            Timestamp: time.Now(),
        },
        PromotionID:   promotionId,
        PromotionType: promotionType,
        Discount:      discount,
        OrderNo:       orderNo,
    }
}

func NewPointsEvent(
    userId int64,
    points int64,
    balance int64,
    source string,
    orderNo string,
    reason string,
    eventType types.MarketingEventType,
) *types.PointsEvent {
    return &types.PointsEvent{
        MarketingEvent: types.MarketingEvent{
            Type:      eventType,
            UserID:    userId,
            Timestamp: time.Now(),
        },
        Points:  points,
        Balance: balance,
        Source:  source,
        OrderNo: orderNo,
        Reason:  reason,
    }
}

func NewMarketingPaymentSuccessEvent(
    userId int64,
    orderNo string,
    paymentNo string,
    amount float64,
    rewardPoints int64,
    couponId int64,
) *types.MarketingPaymentSuccessEvent {
    return &types.MarketingPaymentSuccessEvent{
        MarketingEvent: types.MarketingEvent{
            Type:      types.MarketingPaymentSuccess,
            UserID:    userId,
            Timestamp: time.Now(),
        },
        OrderNo:      orderNo,
        PaymentNo:    paymentNo,
        Amount:       amount,
        RewardPoints: rewardPoints,
        CouponID:     couponId,
    }
}

func NewMarketingPaymentFailedEvent(
    userId int64,
    orderNo string,
    paymentNo string,
    reason string,
) *types.MarketingPaymentFailedEvent {
    return &types.MarketingPaymentFailedEvent{
        MarketingEvent: types.MarketingEvent{
            Type:      types.MarketingPaymentFailed,
            UserID:    userId,
            Timestamp: time.Now(),
        },
        OrderNo:   orderNo,
        PaymentNo: paymentNo,
        Reason:    reason,
    }
}