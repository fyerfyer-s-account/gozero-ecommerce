package logic

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductLogic {
	return &UpdateProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProductLogic) UpdateProduct(in *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
    // Validate input
    if in.Id <= 0 || in.Name == "" || in.Price <= 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    // Check product exists
    existing, err := l.svcCtx.ProductsModel.FindOne(l.ctx, uint64(in.Id))
    if err != nil {
        return nil, zeroerr.ErrProductNotFound
    }

    // Convert images to JSON
    var imagesJSON string
    if len(in.Images) > 0 {
        imagesBytes, err := json.Marshal(in.Images)
        if err != nil {
            return nil, zeroerr.ErrInvalidParam
        }
        imagesJSON = string(imagesBytes)
    }

    // Update product
    existing.Name = in.Name
    existing.Description = sql.NullString{String: in.Description, Valid: in.Description != ""}
    existing.CategoryId = uint64(in.CategoryId)
    existing.Brand = sql.NullString{String: in.Brand, Valid: in.Brand != ""}
    existing.Images = sql.NullString{String: imagesJSON, Valid: len(in.Images) > 0}
    existing.Price = in.Price
    existing.Status = int64(in.Status)

    err = l.svcCtx.ProductsModel.Update(l.ctx, existing)
    if err != nil {
        logx.Errorf("Failed to update product: %v", err)
        return nil, zeroerr.ErrProductUpdateFailed
    }

    return &product.UpdateProductResponse{
        Success: true,
    }, nil
}
