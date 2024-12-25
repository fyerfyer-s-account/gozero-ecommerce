package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return &marketing.ReceiveCouponResponse{}, nil
}
