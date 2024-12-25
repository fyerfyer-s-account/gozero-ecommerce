package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"

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
	// todo: add your logic here and delete this line

	return &marketing.UserCouponsResponse{}, nil
}
