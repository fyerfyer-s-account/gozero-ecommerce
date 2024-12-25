package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"

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
	// todo: add your logic here and delete this line

	return &marketing.VerifyCouponResponse{}, nil
}
