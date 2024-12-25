package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckStockLogic {
	return &CheckStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckStockLogic) CheckStock(in *cart.CheckStockRequest) (*cart.CheckStockResponse, error) {
	// todo: add your logic here and delete this line

	return &cart.CheckStockResponse{}, nil
}
