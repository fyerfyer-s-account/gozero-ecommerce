package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	product "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSkuStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSkuStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSkuStockLogic {
	return &UpdateSkuStockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSkuStockLogic) UpdateSkuStock(req *types.UpdateSkuStockReq) error {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.ProductRpc.UpdateSkuStock(l.ctx, &product.UpdateSkuStockRequest{
		Id:        req.Id,
		Increment: req.Stock,
	})
	return err
}
