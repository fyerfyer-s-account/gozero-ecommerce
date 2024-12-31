package product

import (
	"context"
	"math"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	product "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchProductsLogic {
	return &SearchProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchProductsLogic) SearchProducts(req *types.SearchReq) (*types.SearchResp, error) {
	// Call product RPC
	resp, err := l.svcCtx.ProductRpc.ListProducts(l.ctx, &product.ListProductsRequest{
		CategoryId: req.CategoryId,
		Keyword:    req.Keyword,
		Page:       req.Page,
	})
	if err != nil {
		return nil, err
	}

	// Convert products
	products := make([]types.Product, 0, len(resp.Products))
	for _, p := range resp.Products {
		products = append(products, types.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			CategoryId:  p.CategoryId,
			Brand:       p.Brand,
			Images:      p.Images,
			Price:       p.Price,
			Sales:       int32(p.Sales),
			Status:      int32(p.Status),
			CreatedAt:   p.CreatedAt,
		})
	}

	// Calculate total pages
	totalPages := int32(math.Ceil(float64(resp.Total) / float64(l.svcCtx.Config.PageSize)))

	return &types.SearchResp{
		List:       products,
		Total:      resp.Total,
		Page:       req.Page,
		TotalPages: totalPages,
	}, nil
}
