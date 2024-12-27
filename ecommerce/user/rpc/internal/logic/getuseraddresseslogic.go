package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserAddressesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAddressesLogic {
	return &GetUserAddressesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserAddressesLogic) GetUserAddresses(in *user.GetUserAddressesRequest) (*user.GetUserAddressesResponse, error) {
	// todo: add your logic here and delete this line

	return &user.GetUserAddressesResponse{}, nil
}
