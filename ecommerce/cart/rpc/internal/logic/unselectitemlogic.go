package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnselectItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnselectItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnselectItemLogic {
	return &UnselectItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnselectItemLogic) UnselectItem(in *cart.UnselectItemRequest) (*cart.UnselectItemResponse, error) {
	// todo: add your logic here and delete this line

	return &cart.UnselectItemResponse{}, nil
}
