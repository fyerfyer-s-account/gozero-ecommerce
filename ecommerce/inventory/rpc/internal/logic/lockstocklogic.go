package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type LockStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLockStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LockStockLogic {
	return &LockStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 库存锁定/解锁
func (l *LockStockLogic) LockStock(in *inventory.LockStockRequest) (*inventory.LockStockResponse, error) {
	// todo: add your logic here and delete this line

	return &inventory.LockStockResponse{}, nil
}
