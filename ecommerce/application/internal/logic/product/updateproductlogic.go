package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	product "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductLogic {
	return &UpdateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductLogic) UpdateProduct(req *types.UpdateProductReq) error {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.ProductRpc.UpdateProduct(l.ctx, &product.UpdateProductRequest{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		CategoryId:  req.CategoryId,
		Brand:       req.Brand,
		Images:      req.Images,
		Price:       req.Price,
	})
	return err
}
