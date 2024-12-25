package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelectAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSelectAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectAllLogic {
	return &SelectAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SelectAllLogic) SelectAll(in *cart.SelectAllRequest) (*cart.SelectAllResponse, error) {
	// todo: add your logic here and delete this line

	return &cart.SelectAllResponse{}, nil
}
