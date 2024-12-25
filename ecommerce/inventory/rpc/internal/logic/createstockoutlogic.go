package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateStockOutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateStockOutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateStockOutLogic {
	return &CreateStockOutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateStockOutLogic) CreateStockOut(in *inventory.CreateStockOutRequest) (*inventory.CreateStockOutResponse, error) {
	// todo: add your logic here and delete this line

	return &inventory.CreateStockOutResponse{}, nil
}
