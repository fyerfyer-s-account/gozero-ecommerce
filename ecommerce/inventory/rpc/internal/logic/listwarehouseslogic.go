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
    // Input validation with defaults
    if in.Page <= 0 {
        in.Page = 1
    }
    if in.PageSize <= 0 || in.PageSize > 100 {
        in.PageSize = 20
    }

    // Query warehouses
    warehouses, err := l.svcCtx.WarehousesModel.FindMany(l.ctx, in.Page, in.PageSize)
    if err != nil {
        return nil, err
    }

    // Get total count
    total, err := l.svcCtx.WarehousesModel.Count(l.ctx)
    if err != nil {
        return nil, err
    }

    // Transform to response
    result := make([]*inventory.Warehouse, 0, len(warehouses))
    for _, w := range warehouses {
        result = append(result, &inventory.Warehouse{
            Id:        int64(w.Id),
            Name:      w.Name,
            Address:   w.Address,
            Contact:   w.Contact.String,
            Phone:     w.Phone.String,
            Status:    int32(w.Status),
            CreatedAt: w.CreatedAt.Unix(),
            UpdatedAt: w.UpdatedAt.Unix(),
        })
    }

    return &inventory.ListWarehousesResponse{
        Warehouses: result,
        Total:     total,
    }, nil
}
