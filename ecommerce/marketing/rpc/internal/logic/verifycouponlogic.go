package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCouponLogic {
	return &VerifyCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyCouponLogic) VerifyCoupon(in *marketing.VerifyCouponRequest) (*marketing.VerifyCouponResponse, error) {
    if in.UserId <= 0 || in.CouponId <= 0 || in.Amount <= 0 {
        return nil, zeroerr.ErrInvalidMarketingParam
    }

    // Get coupon
    coupon, err := l.svcCtx.CouponsModel.FindOne(l.ctx, uint64(in.CouponId))
    if err != nil {
        return nil, zeroerr.ErrCouponNotFound
    }

    // Verify coupon status
    now := time.Now()
    if coupon.Status != 1 || 
       (coupon.StartTime.Valid && coupon.StartTime.Time.After(now)) || 
       (coupon.EndTime.Valid && coupon.EndTime.Time.Before(now)) {
        return &marketing.VerifyCouponResponse{
            Valid:   false,
            Message: "Coupon is not valid",
        }, nil
    }

    // Verify user ownership
    userCoupon, err := l.svcCtx.UserCouponsModel.VerifyCoupon(l.ctx, in.UserId, in.CouponId)
    if err != nil || userCoupon == nil {
        return &marketing.VerifyCouponResponse{
            Valid:   false,
            Message: "Coupon not found or already used",
        }, nil
    }

    // Check minimum amount
    if in.Amount < coupon.MinAmount {
        return &marketing.VerifyCouponResponse{
            Valid:   false,
            Message: fmt.Sprintf("Order amount must be greater than %.2f", coupon.MinAmount),
        }, nil
    }

    // Calculate discount
    var discountAmount float64
    switch coupon.Type {
    case 1: // Fixed amount
        discountAmount = coupon.Value
    case 2: // Percentage
        discountAmount = in.Amount * (coupon.Value / 100)
    case 3: // Direct reduction
        discountAmount = coupon.Value
    }

    if discountAmount > in.Amount {
        discountAmount = in.Amount
    }

    return &marketing.VerifyCouponResponse{
        Valid:          true,
        Message:        "Coupon is valid",
        DiscountAmount: discountAmount,
    }, nil
}
