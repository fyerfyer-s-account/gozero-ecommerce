package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCouponLogic {
	return &GetCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCouponLogic) GetCoupon(in *marketing.GetCouponRequest) (*marketing.GetCouponResponse, error) {
    if in.Id <= 0 {
        return nil, zeroerr.ErrInvalidMarketingParam
    }

    coupon, err := l.svcCtx.CouponsModel.FindOne(l.ctx, uint64(in.Id))
    if err != nil {
        l.Logger.Errorf("Failed to get coupon: %v", err)
        return nil, zeroerr.ErrCouponNotFound
    }

    return &marketing.GetCouponResponse{
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
    }, nil
}
