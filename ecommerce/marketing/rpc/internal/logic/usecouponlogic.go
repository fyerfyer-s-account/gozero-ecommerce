package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UseCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUseCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UseCouponLogic {
	return &UseCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UseCouponLogic) UseCoupon(in *marketing.UseCouponRequest) (*marketing.UseCouponResponse, error) {
    if in.UserId <= 0 || in.CouponId <= 0 || in.OrderNo == "" {
        return nil, zeroerr.ErrInvalidMarketingParam
    }

    var success bool
    err := l.svcCtx.UserCouponsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
        // Verify user coupon
        userCoupon, err := l.svcCtx.UserCouponsModel.VerifyCoupon(ctx, in.UserId, in.CouponId)
        if err != nil {
            return zeroerr.ErrCouponNotFound
        }

        // Check if already used
        if userCoupon.Status != 0 {
            return zeroerr.ErrCouponUsed
        }

        // Update user coupon status
        err = l.svcCtx.UserCouponsModel.UpdateStatus(ctx, int64(userCoupon.Id), 1, in.OrderNo)
        if err != nil {
            return err
        }

        // Update coupon used count
        err = l.svcCtx.CouponsModel.IncrUsed(ctx, uint64(in.CouponId))
        if err != nil {
            return err
        }

        success = true
        return nil
    })

    if err != nil {
        l.Logger.Errorf("Failed to use coupon: %v", err)
        return &marketing.UseCouponResponse{Success: false}, err
    }

    // Publish event
    event := types.NewMarketingEvent(types.EventTypeCouponUsed, &types.CouponEventData{
        CouponID: in.CouponId,
        UserID:   in.UserId,
        OrderNo:  in.OrderNo,
    })
    
    if err := l.svcCtx.Producer.PublishCouponEvent(event); err != nil {
        l.Logger.Errorf("Failed to publish coupon used event: %v", err)
    }

    return &marketing.UseCouponResponse{
        Success: success,
    }, nil
}
