package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListStockRecordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListStockRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListStockRecordsLogic {
	return &ListStockRecordsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListStockRecordsLogic) ListStockRecords(in *inventory.ListStockRecordsRequest) (*inventory.ListStockRecordsResponse, error) {
	// todo: add your logic here and delete this line

	return &inventory.ListStockRecordsResponse{}, nil
}
