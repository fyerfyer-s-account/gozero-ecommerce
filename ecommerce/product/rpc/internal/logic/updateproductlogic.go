package logic

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
	if in.Id <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Check if product exists
	_, err := l.svcCtx.ProductsModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, zeroerr.ErrProductNotFound
		}
		return nil, err
	}

	// Prepare updates
	updates := make(map[string]interface{})
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.Description != "" {
		updates["description"] = sql.NullString{String: in.Description, Valid: true}
	}
	if in.CategoryId > 0 {
		updates["category_id"] = in.CategoryId
	}
	if in.Brand != "" {
		updates["brand"] = sql.NullString{String: in.Brand, Valid: true}
	}
	if len(in.Images) > 0 {
		imagesJSON, err := json.Marshal(in.Images)
		if err != nil {
			return nil, err
		}
		updates["images"] = sql.NullString{String: string(imagesJSON), Valid: true}
	}
	if in.Price > 0 {
		updates["price"] = in.Price
	}
	if in.Status > 0 {
		updates["status"] = in.Status
	}

	// Update product
	err = l.svcCtx.ProductsModel.UpdatePartial(l.ctx, uint64(in.Id), updates)
	if err != nil {
		return nil, err
	}

	return &product.UpdateProductResponse{Success: true}, nil
}
