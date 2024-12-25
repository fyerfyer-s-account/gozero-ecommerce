package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetStockLogic {
	return &BatchGetStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchGetStockLogic) BatchGetStock(in *inventory.BatchGetStockRequest) (*inventory.BatchGetStockResponse, error) {
	// todo: add your logic here and delete this line

	return &inventory.BatchGetStockResponse{}, nil
}
