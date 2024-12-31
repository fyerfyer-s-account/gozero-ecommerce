package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	product "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductSkusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductSkusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductSkusLogic {
	return &GetProductSkusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}


func (l *GetProductSkusLogic) GetProductSkus(req *types.GetProductSkusReq) ([]types.Sku, error) {
    // Call product RPC
    skuResp, err := l.svcCtx.ProductRpc.GetSku(l.ctx, &product.GetSkuRequest{
        Id: req.Id,
    })
    if err != nil {
        return nil, err
    }

    // Convert attributes
    attrs := make(map[string]string)
    for _, attr := range skuResp.Sku.Attributes {
        attrs[attr.Key] = attr.Value
    }

    // Convert to API type and return as single-element slice
    return []types.Sku{
        {
            Id:         skuResp.Sku.Id,
            ProductId:  skuResp.Sku.ProductId,
            Name:       skuResp.Sku.SkuCode,
            Code:       skuResp.Sku.SkuCode,
            Price:      skuResp.Sku.Price,
            Stock:      int32(skuResp.Sku.Stock),
            Attributes: attrs,
        },
    }, nil
}
