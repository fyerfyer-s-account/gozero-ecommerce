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
    // Input validation
    if in.Page <= 0 {
        in.Page = 1
    }
    if in.PageSize <= 0 || in.PageSize > 100 {
        in.PageSize = 20
    }

    // Query records with filters
    records, total, err := l.svcCtx.StockRecordsModel.List(
        l.ctx,
        uint64(in.SkuId),
        uint64(in.WarehouseId),
        in.Type,
        in.Page,
        in.PageSize,
    )
    if err != nil {
        return nil, err
    }

    // Transform to response format
    result := make([]*inventory.StockRecord, 0, len(records))
    for _, record := range records {
        result = append(result, &inventory.StockRecord{
            Id:          int64(record.Id),
            SkuId:       int64(record.SkuId),
            WarehouseId: int64(record.WarehouseId),
            Type:        int32(record.Type),
            Quantity:    int32(record.Quantity),
            OrderNo:     record.OrderNo.String,
            Remark:      record.Remark.String,
            CreatedAt:   record.CreatedAt.Unix(),
        })
    }

    return &inventory.ListStockRecordsResponse{
        Records: result,
        Total:   total,
    }, nil
}
