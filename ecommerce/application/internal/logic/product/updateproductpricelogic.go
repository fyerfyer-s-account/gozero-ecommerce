package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	product "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductPriceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductPriceLogic {
	return &UpdateProductPriceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductPriceLogic) UpdateProductPrice(req *types.UpdateProductPriceReq) error {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.ProductRpc.UpdateProductPrice(l.ctx, &product.UpdateProductPriceRequest{
		Id:    req.Id,
		Price: req.Price,
	})
	return err
}
