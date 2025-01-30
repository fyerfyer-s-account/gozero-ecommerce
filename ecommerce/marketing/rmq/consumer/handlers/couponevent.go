package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CouponEventHandler struct {
    logger            *zerolog.Logger
    couponsModel      model.CouponsModel
    userCouponsModel  model.UserCouponsModel
}

func NewCouponEventHandler(
    couponsModel model.CouponsModel,
    userCouponsModel model.UserCouponsModel,
) *CouponEventHandler {
    return &CouponEventHandler{
        logger:           zerolog.GetLogger(),
        couponsModel:     couponsModel,
        userCouponsModel: userCouponsModel,
    }
}

func (h *CouponEventHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.CouponEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return fmt.Errorf("failed to unmarshal coupon event: %w", err)
    }

    fields := map[string]interface{}{
        "type":        event.Type,
        "user_id":     event.UserID,
        "coupon_id":   event.CouponID,
        "coupon_code": event.CouponCode,
    }
    h.logger.Info(ctx, "Processing coupon event", fields)

    switch event.Type {
    case types.CouponIssued:
        return h.handleCouponIssued(ctx, &event)
    case types.CouponUsed:
        return h.handleCouponUsed(ctx, &event)
    case types.CouponExpired:
        return h.handleCouponExpired(ctx, &event)
    case types.CouponCancelled:
        return h.handleCouponCancelled(ctx, &event)
    default:
        return fmt.Errorf("unknown coupon event type: %s", event.Type)
    }
}

func (h *CouponEventHandler) handleCouponIssued(ctx context.Context, event *types.CouponEvent) error {
    // Check if user has already received this coupon
    count, err := h.userCouponsModel.CountUserCoupon(ctx, event.UserID, event.CouponID)
    if err != nil {
        return fmt.Errorf("failed to check user coupon count: %w", err)
    }

    coupon, err := h.couponsModel.FindOne(ctx, uint64(event.CouponID))
    if err != nil {
        return fmt.Errorf("failed to get coupon: %w", err)
    }

    // Check if user has exceeded limit
    if coupon.PerLimit > 0 && count >= coupon.PerLimit {
        return fmt.Errorf("user has reached coupon limit")
    }

    // Create user coupon record
    userCoupon := &model.UserCoupons{
        UserId:    uint64(event.UserID),
        CouponId:  uint64(event.CouponID),
        Status:    0, // Unused
        CreatedAt: time.Now(),
    }

    err = h.userCouponsModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        // Insert user coupon
        if _, err := h.userCouponsModel.Insert(ctx, userCoupon); err != nil {
            return err
        }

        // Increment received count
        if err := h.couponsModel.IncrReceived(ctx, uint64(event.CouponID)); err != nil {
            return err
        }

        return nil
    })

    return err
}

func (h *CouponEventHandler) handleCouponUsed(ctx context.Context, event *types.CouponEvent) error {
    // Verify coupon exists and is unused
    userCoupon, err := h.userCouponsModel.VerifyCoupon(ctx, event.UserID, event.CouponID)
    if err != nil {
        return fmt.Errorf("failed to verify coupon: %w", err)
    }

    // Update coupon status to used
    err = h.userCouponsModel.UpdateStatus(ctx, int64(userCoupon.Id), 1, event.OrderNo)
    if err != nil {
        return fmt.Errorf("failed to update coupon status: %w", err)
    }

    // Increment used count
    if err := h.couponsModel.IncrUsed(ctx, uint64(event.CouponID)); err != nil {
        return fmt.Errorf("failed to increment used count: %w", err)
    }

    return nil
}

func (h *CouponEventHandler) handleCouponExpired(ctx context.Context, event *types.CouponEvent) error {
    userCoupon, err := h.userCouponsModel.VerifyCoupon(ctx, event.UserID, event.CouponID)
    if err != nil {
        return fmt.Errorf("failed to verify coupon: %w", err)
    }

    // Update coupon status to expired
    err = h.userCouponsModel.UpdateStatus(ctx, int64(userCoupon.Id), 2, "")
    if err != nil {
        return fmt.Errorf("failed to update coupon status: %w", err)
    }

    return nil
}

func (h *CouponEventHandler) handleCouponCancelled(ctx context.Context, event *types.CouponEvent) error {
    userCoupon, err := h.userCouponsModel.VerifyCoupon(ctx, event.UserID, event.CouponID)
    if err != nil {
        return fmt.Errorf("failed to verify coupon: %w", err)
    }

    err = h.couponsModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        // Update coupon status to cancelled
        if err := h.userCouponsModel.UpdateStatus(ctx, int64(userCoupon.Id), 3, ""); err != nil {
            return err
        }

        // Decrement received count
        if err := h.couponsModel.DecrReceived(ctx, uint64(event.CouponID)); err != nil {
            return err
        }

        return nil
    })

    return err
}