package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 商品管理
func (l *CreateProductLogic) CreateProduct(in *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	// Validate input
	if in.Name == "" || in.CategoryId <= 0 || in.Price <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Check category exists
	_, err := l.svcCtx.CategoriesModel.FindOne(l.ctx, uint64(in.CategoryId))
	if err != nil {
		return nil, zeroerr.ErrCategoryNotFound
	}

	// Check product name uniqueness
	existingProd, err := l.svcCtx.ProductsModel.FindOneByName(l.ctx, in.Name)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}
	if existingProd != nil {
		return nil, zeroerr.ErrProductDuplicate
	}

	// Create product
	var imagesJSON string
	if len(in.Images) > 0 {
		imagesBytes, err := json.Marshal(in.Images)
		if err != nil {
			return nil, zeroerr.ErrInvalidParam
		}
		imagesJSON = string(imagesBytes)
	}
	p := &model.Products{
		Name:        in.Name,
		Brief:       sql.NullString{String: in.Brief, Valid: in.Brief != ""},
		Description: sql.NullString{String: in.Description, Valid: in.Description != ""},
		CategoryId:  uint64(in.CategoryId),
		Brand:       sql.NullString{String: in.Brand, Valid: in.Brand != ""},
		Images:      sql.NullString{String: imagesJSON, Valid: len(in.Images) > 0},
		Price:       in.Price,
		Status:      1, // Default to active
	}

	result, err := l.svcCtx.ProductsModel.Insert(l.ctx, p)
	if err != nil {
		logx.Errorf("Failed to create product: %v", err)
		return nil, zeroerr.ErrProductCreateFailed
	}

	productId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Create SKUs if attributes provided
	if len(in.SkuAttributes) > 0 {
		skus := generateSkus(uint64(productId), in.Price, in.SkuAttributes)
		err = l.svcCtx.SkusModel.BatchInsert(l.ctx, skus)
		if err != nil {
			// Cleanup on SKU creation failure
			_ = l.svcCtx.ProductsModel.Delete(l.ctx, uint64(productId))
			return nil, zeroerr.ErrSkuCreateFailed
		}
	}

	return &product.CreateProductResponse{
		Id: productId,
	}, nil
}

func generateSkus(productId uint64, basePrice float64, attrs []*product.SkuAttribute) []*model.Skus {
	skus := make([]*model.Skus, 1)

	// Single SKU for now
	attrJSON, _ := json.Marshal(attrs)
	skus[0] = &model.Skus{
		ProductId:  productId,
		SkuCode:    fmt.Sprintf("SKU-%d-1", productId),
		Price:      basePrice,
		Stock:      0,
		Attributes: string(attrJSON),
	}

	return skus
}
