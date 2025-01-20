package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCouponsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponsLogic {
	return &UserCouponsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserCouponsLogic) UserCoupons(in *marketing.UserCouponsRequest) (*marketing.UserCouponsResponse, error) {
    if in.UserId <= 0 {
        return nil, zeroerr.ErrInvalidMarketingParam
    }

    if in.Page <= 0 {
        in.Page = 1
    }
    if in.PageSize <= 0 || in.PageSize > 100 {
        in.PageSize = 10
    }

    // Get user coupons
    userCoupons, err := l.svcCtx.UserCouponsModel.FindByUserAndStatus(l.ctx, in.UserId, in.Status, in.Page, in.PageSize)
    if err != nil {
        l.Logger.Errorf("Failed to get user coupons: %v", err)
        return nil, err
    }

    // Get total count
    total, err := l.svcCtx.UserCouponsModel.CountByUser(l.ctx, in.UserId, in.Status)
    if err != nil {
        l.Logger.Errorf("Failed to get user coupons count: %v", err)
        return nil, err
    }

    // Get coupon details
    var couponIds []uint64
    for _, uc := range userCoupons {
        couponIds = append(couponIds, uc.CouponId)
    }

    couponsMap, err := l.svcCtx.CouponsModel.FindManyByIds(l.ctx, couponIds)
    if err != nil {
        l.Logger.Errorf("Failed to get coupons details: %v", err)
        return nil, err
    }

    // Build response
    var result []*marketing.UserCoupon
    for _, uc := range userCoupons {
        coupon := couponsMap[uc.CouponId]
        if coupon == nil {
            continue
        }

        result = append(result, &marketing.UserCoupon{
            Id:        int64(uc.Id),
            UserId:    int64(uc.UserId),
            CouponId:  int64(uc.CouponId),
            Status:    int32(uc.Status),
            UsedTime:  uc.UsedTime.Time.Unix(),
            CreatedAt: uc.CreatedAt.Unix(),
            Coupon: &marketing.Coupon{
                Id:        int64(coupon.Id),
                Name:      coupon.Name,
                Code:      coupon.Code,
                Type:      int32(coupon.Type),
                Value:     coupon.Value,
                MinAmount: coupon.MinAmount,
                Status:    int32(coupon.Status),
                StartTime: coupon.StartTime.Time.Unix(),
                EndTime:   coupon.EndTime.Time.Unix(),
                Total:     int32(coupon.Total),
                Received:  int32(coupon.Received),
                Used:      int32(coupon.Used),
                PerLimit:  coupon.PerLimit,
                UserLimit: int32(coupon.UserLimit),
                CreatedAt: coupon.CreatedAt.Unix(),
                UpdatedAt: coupon.UpdatedAt.Unix(),
            },
        })
    }

    return &marketing.UserCouponsResponse{
        Coupons: result,
        Total:   total,
    }, nil
}
