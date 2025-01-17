package logic

import (
	"context"
	"database/sql"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
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
    // Input validation
    if in.Id <= 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    // Check if warehouse exists
    warehouse, err := l.svcCtx.WarehousesModel.FindOne(l.ctx, uint64(in.Id))
    if err != nil {
        if err == model.ErrNotFound {
            return nil, zeroerr.ErrNotFound
        }
        return nil, err
    }

    // Update fields if provided
    if len(in.Name) > 0 {
        warehouse.Name = in.Name
    }
    if len(in.Address) > 0 {
        warehouse.Address = in.Address
    }
    if len(in.Contact) > 0 {
        warehouse.Contact = sql.NullString{String: in.Contact, Valid: true}
    }
    if len(in.Phone) > 0 {
        warehouse.Phone = sql.NullString{String: in.Phone, Valid: true}
    }
    if in.Status > 0 {
        warehouse.Status = int64(in.Status)
    }

    // Update warehouse
    err = l.svcCtx.WarehousesModel.Update(l.ctx, warehouse)
    if err != nil {
        return nil, err
    }

    return &inventory.UpdateWarehouseResponse{
        Success: true,
    }, nil
}