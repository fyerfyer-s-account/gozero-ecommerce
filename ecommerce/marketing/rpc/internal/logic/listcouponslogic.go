package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCouponsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCouponsLogic {
	return &ListCouponsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCouponsLogic) ListCoupons(in *marketing.ListCouponsRequest) (*marketing.ListCouponsResponse, error) {
    if in.Page <= 0 {
        in.Page = 1
    }
    if in.PageSize <= 0 || in.PageSize > 100 {
        in.PageSize = 10
    }

    coupons, err := l.svcCtx.CouponsModel.FindMany(l.ctx, in.Status, in.Page, in.PageSize)
    if err != nil {
        l.Logger.Errorf("Failed to get coupons list: %v", err)
        return nil, zeroerr.ErrCouponNotFound
    }

    total, err := l.svcCtx.CouponsModel.Count(l.ctx, in.Status)
    if err != nil {
        l.Logger.Errorf("Failed to get coupons count: %v", err)
        return nil, zeroerr.ErrCouponNotFound
    }

    var result []*marketing.Coupon
    for _, c := range coupons {
        result = append(result, &marketing.Coupon{
            Id:        int64(c.Id),
            Name:      c.Name,
            Code:      c.Code,
            Type:      int32(c.Type),
            Value:     c.Value,
            MinAmount: c.MinAmount,
            Status:    int32(c.Status),
            StartTime: c.StartTime.Time.Unix(),
            EndTime:   c.EndTime.Time.Unix(),
            Total:     int32(c.Total),
            Received:  int32(c.Received),
            Used:      int32(c.Used),
            PerLimit:  c.PerLimit,
            UserLimit: int32(c.UserLimit),
            CreatedAt: c.CreatedAt.Unix(),
            UpdatedAt: c.UpdatedAt.Unix(),
        })
    }

    return &marketing.ListCouponsResponse{
        Coupons: result,
        Total:   total,
    }, nil
}
