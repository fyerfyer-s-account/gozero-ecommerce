package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnlockStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnlockStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlockStockLogic {
	return &UnlockStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnlockStockLogic) UnlockStock(in *inventory.UnlockStockRequest) (*inventory.UnlockStockResponse, error) {
	// todo: add your logic here and delete this line

	return &inventory.UnlockStockResponse{}, nil
}
