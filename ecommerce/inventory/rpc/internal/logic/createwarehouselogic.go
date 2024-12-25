package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWarehouseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWarehouseLogic {
	return &CreateWarehouseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 仓储管理
func (l *CreateWarehouseLogic) CreateWarehouse(in *inventory.CreateWarehouseRequest) (*inventory.CreateWarehouseResponse, error) {
	// todo: add your logic here and delete this line

	return &inventory.CreateWarehouseResponse{}, nil
}
