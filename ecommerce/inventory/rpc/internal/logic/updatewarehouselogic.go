package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWarehouseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWarehouseLogic {
	return &UpdateWarehouseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateWarehouseLogic) UpdateWarehouse(in *inventory.UpdateWarehouseRequest) (*inventory.UpdateWarehouseResponse, error) {
	// todo: add your logic here and delete this line

	return &inventory.UpdateWarehouseResponse{}, nil
}
