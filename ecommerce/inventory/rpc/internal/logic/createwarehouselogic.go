package logic

import (
	"context"
	"database/sql"
	"strings"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

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
    // Input validation
    if len(strings.TrimSpace(in.Name)) == 0 || len(strings.TrimSpace(in.Address)) == 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    // Check if warehouse name already exists
    existingWarehouse, err := l.svcCtx.WarehousesModel.FindByName(l.ctx, in.Name)
    if err != nil && err != model.ErrNotFound {
        return nil, err
    }
    if existingWarehouse != nil {
        return nil, zeroerr.ErrDuplicateWarehouse
    }

    // Create warehouse
    warehouse := &model.Warehouses{
        Name:    in.Name,
        Address: in.Address,
        Contact: sql.NullString{
            String: in.Contact,
            Valid:  len(in.Contact) > 0,
        },
        Phone: sql.NullString{
            String: in.Phone,
            Valid:  len(in.Phone) > 0,
        },
        Status: 1, // Default to active
    }

    result, err := l.svcCtx.WarehousesModel.Insert(l.ctx, warehouse)
    if err != nil {
        return nil, zeroerr.ErrWarehouseCreateFailed
    }

    id, err := result.LastInsertId()
    if err != nil {
        return nil, zeroerr.ErrWarehouseCreateFailed
    }

    return &inventory.CreateWarehouseResponse{
        Id: id,
    }, nil
}
