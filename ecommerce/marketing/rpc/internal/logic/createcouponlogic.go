package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCouponLogic {
	return &CreateCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 优惠券管理
func (l *CreateCouponLogic) CreateCoupon(in *marketing.CreateCouponRequest) (*marketing.CreateCouponResponse, error) {
	// todo: add your logic here and delete this line

	return &marketing.CreateCouponResponse{}, nil
}
