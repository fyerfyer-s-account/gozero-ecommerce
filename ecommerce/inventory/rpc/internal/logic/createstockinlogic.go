package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateStockInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateStockInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateStockInLogic {
	return &CreateStockInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 入库/出库
func (l *CreateStockInLogic) CreateStockIn(in *inventory.CreateStockInRequest) (*inventory.CreateStockInResponse, error) {
	// todo: add your logic here and delete this line

	return &inventory.CreateStockInResponse{}, nil
}
