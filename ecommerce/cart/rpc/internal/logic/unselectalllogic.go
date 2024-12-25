package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnselectAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnselectAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnselectAllLogic {
	return &UnselectAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnselectAllLogic) UnselectAll(in *cart.UnselectAllRequest) (*cart.UnselectAllResponse, error) {
	// todo: add your logic here and delete this line

	return &cart.UnselectAllResponse{}, nil
}
