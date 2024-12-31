package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	product "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductLogic) GetProduct(req *types.GetProductReq) (*types.Product, error) {
	resp, err := l.svcCtx.ProductRpc.GetProduct(l.ctx, &product.GetProductRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &types.Product{
		Id:          resp.Product.Id,
		Name:        resp.Product.Name,
		Description: resp.Product.Description,
		CategoryId:  resp.Product.CategoryId,
		Brand:       resp.Product.Brand,
		Images:      resp.Product.Images,
		Price:       resp.Product.Price,
		Stock:       calculateTotalStock(resp.Skus),
		Sales:       int32(resp.Product.Sales),
		Status:      int32(resp.Product.Status),
		CreatedAt:   resp.Product.CreatedAt,
	}, nil
}

func calculateTotalStock(skus []*product.Sku) int32 {
	var total int32
	for _, sku := range skus {
		total += int32(sku.Stock)
	}
	return total
}
