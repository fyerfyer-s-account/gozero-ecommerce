package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type PaymentSuccessHandler struct {
    logger            *zerolog.Logger
    userPointsModel   model.UserPointsModel
    userCouponsModel  model.UserCouponsModel
    pointsRecordModel model.PointsRecordsModel
}

func NewPaymentSuccessHandler(
    userPointsModel model.UserPointsModel,
    userCouponsModel model.UserCouponsModel,
    pointsRecordModel model.PointsRecordsModel,
) *PaymentSuccessHandler {
    return &PaymentSuccessHandler{
        logger:            zerolog.GetLogger(),
        userPointsModel:   userPointsModel,
        userCouponsModel:  userCouponsModel,
        pointsRecordModel: pointsRecordModel,
    }
}

func (h *PaymentSuccessHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.MarketingPaymentSuccessEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return fmt.Errorf("failed to unmarshal payment success event: %w", err)
    }

    fields := map[string]interface{}{
        "user_id":       event.UserID,
        "order_no":      event.OrderNo,
        "payment_no":    event.PaymentNo,
        "amount":        event.Amount,
        "reward_points": event.RewardPoints,
        "coupon_id":     event.CouponID,
    }
    h.logger.Info(ctx, "Processing payment success event", fields)

    // Handle in transaction
    err := h.userPointsModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        // 1. Handle reward points if any
        if event.RewardPoints > 0 {
            if err := h.handleRewardPoints(ctx, &event); err != nil {
                return fmt.Errorf("failed to handle reward points: %w", err)
            }
        }

        // 2. Update coupon status if used
        if event.CouponID > 0 {
            if err := h.handleCouponUsage(ctx, &event); err != nil {
                return fmt.Errorf("failed to handle coupon usage: %w", err)
            }
        }

        return nil
    })

    if err != nil {
        h.logger.Error(ctx, "Failed to process payment success event", err, fields)
        return err
    }

    h.logger.Info(ctx, "Successfully processed payment success event", fields)
    return nil
}

func (h *PaymentSuccessHandler) handleRewardPoints(ctx context.Context, event *types.MarketingPaymentSuccessEvent) error {
    // Initialize user points if not exists
    _, err := h.userPointsModel.FindOne(ctx, uint64(event.UserID))
    if err == model.ErrNotFound {
        if err := h.userPointsModel.InitUserPoints(ctx, event.UserID); err != nil {
            return fmt.Errorf("failed to init user points: %w", err)
        }
    } else if err != nil {
        return fmt.Errorf("failed to check user points: %w", err)
    }

    // Create points record
    record := &model.PointsRecords{
        UserId:    uint64(event.UserID),
        Points:    event.RewardPoints,
        Type:      1, // Points earned
        Source:    "payment_reward",
        OrderNo:   sql.NullString{String: event.OrderNo, Valid: true},
        Remark:    sql.NullString{String: fmt.Sprintf("Payment reward for order %s", event.OrderNo), Valid: true},
        CreatedAt: time.Now(),
    }

    if _, err := h.pointsRecordModel.Insert(ctx, record); err != nil {
        return fmt.Errorf("failed to create points record: %w", err)
    }

    // Increment user points
    if err := h.userPointsModel.IncrPoints(ctx, event.UserID, event.RewardPoints); err != nil {
        return fmt.Errorf("failed to increment points: %w", err)
    }

    return nil
}

func (h *PaymentSuccessHandler) handleCouponUsage(ctx context.Context, event *types.MarketingPaymentSuccessEvent) error {
    // Get user coupon
    coupon, err := h.userCouponsModel.VerifyCoupon(ctx, event.UserID, event.CouponID)
    if err != nil {
        return fmt.Errorf("failed to verify coupon: %w", err)
    }

    // Update coupon status to used
    if err := h.userCouponsModel.UpdateStatus(ctx, int64(coupon.Id), 1, event.OrderNo); err != nil {
        return fmt.Errorf("failed to update coupon status: %w", err)
    }

    return nil
}