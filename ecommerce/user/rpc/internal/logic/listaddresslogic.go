package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAddressLogic {
	return &ListAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListAddressLogic) ListAddress(in *user.ListAddressRequest) (*user.ListAddressResponse, error) {
	// todo: add your logic here and delete this line

	return &user.ListAddressResponse{}, nil
}
