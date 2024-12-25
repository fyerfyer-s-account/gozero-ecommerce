package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWarehousesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListWarehousesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWarehousesLogic {
	return &ListWarehousesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListWarehousesLogic) ListWarehouses(in *inventory.ListWarehousesRequest) (*inventory.ListWarehousesResponse, error) {
	// todo: add your logic here and delete this line

	return &inventory.ListWarehousesResponse{}, nil
}
