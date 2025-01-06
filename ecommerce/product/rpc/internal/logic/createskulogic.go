package logic

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSkuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSkuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSkuLogic {
	return &CreateSkuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SKU管理
func (l *CreateSkuLogic) CreateSku(in *product.CreateSkuRequest) (*product.CreateSkuResponse, error) {
	// Validate input
	if in.ProductId <= 0 || in.SkuCode == "" || in.Price <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Check product exists
	_, err := l.svcCtx.ProductsModel.FindOne(l.ctx, uint64(in.ProductId))
	if err != nil {
		return nil, zeroerr.ErrProductNotFound
	}

	// Check SKU code uniqueness
	existingSku, err := l.svcCtx.SkusModel.FindOneBySkuCode(l.ctx, in.SkuCode)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}
	if existingSku != nil {
		return nil, zeroerr.ErrSkuDuplicate
	}

	// Convert attributes to JSON
	var attrsJSON string
	if len(in.Attributes) > 0 {
		attrBytes, err := json.Marshal(in.Attributes)
		if err != nil {
			return nil, zeroerr.ErrInvalidAttributes
		}
		attrsJSON = string(attrBytes)
	}

	// Create SKU
	sku := &model.Skus{
		ProductId:  uint64(in.ProductId),
		SkuCode:    in.SkuCode,
		Price:      in.Price,
		Stock:      in.Stock,
		Attributes: sql.NullString {
			String: string(attrsJSON),
			Valid:  true,
		},
		Sales:      0,
	}

	result, err := l.svcCtx.SkusModel.Insert(l.ctx, sku)
	if err != nil {
		logx.Errorf("Failed to create SKU: %v", err)
		return nil, zeroerr.ErrSkuCreateFailed
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &product.CreateSkuResponse{
		Id: id,
	}, nil
}
