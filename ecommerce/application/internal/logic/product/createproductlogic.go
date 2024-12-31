package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	product "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProductLogic) CreateProduct(req *types.CreateProductReq) (resp *types.CreateProductResp, err error) {
	// todo: add your logic here and delete this line
	// Convert attributes
	attrs := make([]*product.SkuAttribute, 0, len(req.Attributes))
	for _, attr := range req.Attributes {
		attrs = append(attrs, &product.SkuAttribute{
			Key:   attr.Key,
			Value: attr.Value,
		})
	}

	// Call RPC
	res, err := l.svcCtx.ProductRpc.CreateProduct(l.ctx, &product.CreateProductRequest{
		Name:          req.Name,
		Description:   req.Description,
		CategoryId:    req.CategoryId,
		Brand:         req.Brand,
		Images:        req.Images,
		Price:         req.Price,
		SkuAttributes: attrs,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateProductResp{
		Id: res.Id,
	}, nil
}
