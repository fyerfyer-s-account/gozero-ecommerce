package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	product "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSkuPriceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSkuPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSkuPriceLogic {
	return &UpdateSkuPriceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSkuPriceLogic) UpdateSkuPrice(req *types.UpdateSkuPriceReq) error {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.ProductRpc.UpdateSkuPrice(l.ctx, &product.UpdateSkuPriceRequest{
		Id:    req.Id,
		Price: req.Price,
	})
	return err
}
