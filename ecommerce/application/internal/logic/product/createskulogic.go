package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	product "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSkuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSkuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSkuLogic {
	return &CreateSkuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSkuLogic) CreateSku(req *types.CreateSkuReq) (resp *types.CreateSkuResp, err error) {
	// todo: add your logic here and delete this line
	attrs := make([]*product.SkuAttribute, 0, len(req.Attributes))
	for _, attr := range req.Attributes {
		attrs = append(attrs, &product.SkuAttribute{
			Key:   attr.Key,
			Value: attr.Value,
		})
	}

	// Call RPC
	res, err := l.svcCtx.ProductRpc.CreateSku(l.ctx, &product.CreateSkuRequest{
		ProductId:  req.ProductId,
		SkuCode:    req.SkuCode,
		Price:      req.Price,
		Stock:      req.Stock,
		Attributes: attrs,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateSkuResp{
		Id: res.Id,
	}, nil
}
