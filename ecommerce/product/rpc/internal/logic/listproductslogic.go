package logic

import (
	"context"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductsLogic {
	return &ListProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListProductsLogic) ListProducts(in *product.ListProductsRequest) (*product.ListProductsResponse, error) {
	if in.PageSize <= 0 {
		in.PageSize = 10
	}
	if in.Page <= 0 {
		in.Page = 1
	}

	var products []*model.Products
	var err error

	// Get total count
	total, err := l.svcCtx.ProductsModel.Count(l.ctx, uint64(in.CategoryId), in.Keyword)
	if err != nil {
		logx.Errorf("Failed to get products count: %v", err)
		return nil, err
	}

	// Get products with filters
	if in.CategoryId > 0 {
		products, err = l.svcCtx.ProductsModel.FindManyByCategoryId(l.ctx, uint64(in.CategoryId), int(in.Page), int(in.PageSize))
	} else if in.Keyword != "" {
		products, err = l.svcCtx.ProductsModel.SearchByKeyword(l.ctx, in.Keyword, int(in.Page), int(in.PageSize))
	} else {
		products, err = l.svcCtx.ProductsModel.GeneralSearch(l.ctx, int(in.Page), int(in.PageSize))
	}

	if err != nil {
		logx.Errorf("Failed to get products: %v", err)
		return nil, zeroerr.ErrProductNotFound
	}

	// Convert to proto messages
	pbProducts := make([]*product.Product, 0, len(products))
	for _, p := range products {
		pbProduct := &product.Product{
			Id:          int64(p.Id),
			Name:        p.Name,
			Description: p.Description.String,
			CategoryId:  int64(p.CategoryId),
			Brand:       p.Brand.String,
			Price:       p.Price,
			Sales:       p.Sales,
			Status:      p.Status,
			CreatedAt:   p.CreatedAt.Unix(),
			UpdatedAt:   p.UpdatedAt.Unix(),
		}
		if p.Images.Valid {
			var images []string
			if err := json.Unmarshal([]byte(p.Images.String), &images); err == nil {
				pbProduct.Images = images
			}
		}
		pbProducts = append(pbProducts, pbProduct)
	}

	return &product.ListProductsResponse{
		Total:    total,
		Products: pbProducts,
	}, nil
}
