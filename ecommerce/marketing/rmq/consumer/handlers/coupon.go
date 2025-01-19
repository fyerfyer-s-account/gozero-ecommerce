package handlers

import (
	"context"
	"database/sql"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CouponHandler struct {
    couponsModel     model.CouponsModel
    userCouponsModel model.UserCouponsModel
}

func NewCouponHandler(couponsModel model.CouponsModel, userCouponsModel model.UserCouponsModel) *CouponHandler {
    return &CouponHandler{
        couponsModel:     couponsModel,
        userCouponsModel: userCouponsModel,
    }
}

func (h *CouponHandler) HandleCouponReceived(event *types.MarketingEvent) error {
    data, ok := event.Data.(*types.CouponEventData)
    if !ok {
        return zeroerr.ErrInvalidEventData
    }

    return h.couponsModel.Trans(context.Background(), func(ctx context.Context, session sqlx.Session) error {
        // Check coupon exists and available
        coupon, err := h.couponsModel.FindOne(ctx, uint64(data.CouponID))
        if err != nil {
            return err
        }
        if coupon.Status != 1 || coupon.Received >= coupon.Total {
            return zeroerr.ErrCouponUnavailable
        }

        // Create user coupon record
        _, err = h.userCouponsModel.Insert(ctx, &model.UserCoupons{
            UserId:   uint64(data.UserID),
            CouponId: uint64(data.CouponID),
            Status:   0, // unused
        })
        if err != nil {
            return err
        }

        // Update coupon received count
        coupon.Received++
        return h.couponsModel.Update(ctx, coupon)
    })
}

func (h *CouponHandler) HandleCouponUsed(event *types.MarketingEvent) error {
    data, ok := event.Data.(*types.CouponEventData)
    if !ok {
        return zeroerr.ErrInvalidEventData
    }

    return h.couponsModel.Trans(context.Background(), func(ctx context.Context, session sqlx.Session) error {
        // Update user coupon status
        userCoupon, err := h.userCouponsModel.FindOne(ctx, uint64(data.CouponID))
        if err != nil {
            return err
        }

        userCoupon.Status = 1 // used
        userCoupon.UsedTime = sql.NullTime{Time: time.Now(), Valid: true}
        userCoupon.OrderNo = sql.NullString{String: data.OrderNo, Valid: true}
        
        if err := h.userCouponsModel.Update(ctx, userCoupon); err != nil {
            return err
        }

        // Update coupon usage count
        coupon, err := h.couponsModel.FindOne(ctx, uint64(data.CouponID))
        if err != nil {
            return err
        }

        coupon.Used++
        return h.couponsModel.Update(ctx, coupon)
    })
}