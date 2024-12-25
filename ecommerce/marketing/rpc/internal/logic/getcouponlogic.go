package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"

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
	// todo: add your logic here and delete this line

	return &marketing.GetCouponResponse{}, nil
}
