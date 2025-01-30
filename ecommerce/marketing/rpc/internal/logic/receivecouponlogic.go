package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ReceiveCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReceiveCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveCouponLogic {
	return &ReceiveCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReceiveCouponLogic) ReceiveCoupon(in *marketing.ReceiveCouponRequest) (*marketing.ReceiveCouponResponse, error) {
    if in.UserId <= 0 || in.CouponId <= 0 {
        return nil, zeroerr.ErrInvalidMarketingParam
    }

    var couponCode string
    err := l.svcCtx.CouponsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
        // Get coupon
        coupon, err := l.svcCtx.CouponsModel.FindOne(ctx, uint64(in.CouponId))
        if err != nil {
            return zeroerr.ErrCouponNotFound
        }
        couponCode = coupon.Code

        // Check coupon status and time
        now := time.Now()
        if coupon.Status != 1 || 
           (coupon.StartTime.Valid && coupon.StartTime.Time.After(now)) || 
           (coupon.EndTime.Valid && coupon.EndTime.Time.Before(now)) {
            return zeroerr.ErrCouponUnavailable
        }

        // Check quantity limit
        if coupon.Received >= coupon.Total {
            return zeroerr.ErrCouponUnavailable
        }

        // Check user limit
        if coupon.PerLimit > 0 {
            count, err := l.svcCtx.UserCouponsModel.CountUserCoupon(ctx, in.UserId, in.CouponId)
            if err != nil {
                return err
            }
            if count >= coupon.UserLimit {
                return zeroerr.ErrExceedCouponLimit
            }
        }

        // Create user coupon record
        _, err = l.svcCtx.UserCouponsModel.Insert(ctx, &model.UserCoupons{
            UserId:    uint64(in.UserId),
            CouponId:  uint64(in.CouponId),
            Status:    0,
            UsedTime:  sql.NullTime{},
            OrderNo:   sql.NullString{},
            CreatedAt: time.Now(),
        })
        if err != nil {
            return err
        }

        // Update coupon received count
        if err := l.svcCtx.CouponsModel.IncrReceived(ctx, uint64(in.CouponId)); err != nil {
            return err
        }

        // Publish coupon issued event
        couponEvent := &types.CouponEvent{
            MarketingEvent: types.MarketingEvent{
                Type:      types.CouponIssued,
                UserID:    in.UserId,
                Timestamp: time.Now(),
            },
            CouponID:   in.CouponId,
            CouponCode: couponCode,
            Amount:     coupon.Value,
        }

        if err := l.svcCtx.Producer.PublishCouponEvent(ctx, couponEvent); err != nil {
            logx.Errorf("Failed to publish coupon issued event: %v", err)
            // Don't return error as this shouldn't affect the main transaction
        }

        return nil
    })

    if err != nil {
        l.Logger.Errorf("Failed to receive coupon: %v", err)
        return &marketing.ReceiveCouponResponse{Success: false}, err
    }

    return &marketing.ReceiveCouponResponse{
        Success: true,
    }, nil
}