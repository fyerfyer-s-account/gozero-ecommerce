package logic

import (
	"context"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSkuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSkuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSkuLogic {
	return &UpdateSkuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateSkuLogic) UpdateSku(in *product.UpdateSkuRequest) (*product.UpdateSkuResponse, error) {
	if in.Id <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Check if SKU exists
	_, err := l.svcCtx.SkusModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, zeroerr.ErrSkuNotFound
	}

	updates := make(map[string]interface{})
	if in.Price > 0 {
		updates["price"] = in.Price
	}
	if in.StockIncrement != 0 {
		updates["stock"] = in.StockIncrement
	}
	if in.SalesIncrement != 0 {
		updates["sales"] = in.SalesIncrement
	}

	err = l.svcCtx.SkusModel.UpdateSkus(l.ctx, uint64(in.Id), updates)
	if err != nil {
		return nil, err
	}

	return &product.UpdateSkuResponse{Success: true}, nil
}
