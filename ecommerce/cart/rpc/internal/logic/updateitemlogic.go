package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateItemLogic {
	return &UpdateItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateItemLogic) UpdateItem(in *cart.UpdateItemRequest) (*cart.UpdateItemResponse, error) {
    if in.UserId <= 0 || in.ProductId <= 0 || in.SkuId <= 0 || in.Quantity <= 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    // Check if item exists
    _, err := l.svcCtx.CartItemsModel.FindOneByUserIdSkuId(l.ctx, uint64(in.UserId), uint64(in.SkuId))
    if err != nil {
        return nil, zeroerr.ErrItemNotFound
    }

    // Check stock
    sku, err := l.svcCtx.ProductRpc.GetSku(l.ctx, &product.GetSkuRequest{
        Id: in.SkuId,
    })
    if err != nil {
        return nil, err
    }

    if sku.Sku.Stock < int64(in.Quantity) {
        return nil, zeroerr.ErrItemOutOfStock
    }

    // Update quantity
    err = l.svcCtx.CartItemsModel.UpdateQuantity(l.ctx, uint64(in.UserId), uint64(in.ProductId), uint64(in.SkuId), int64(in.Quantity))
    if err != nil {
        return nil, zeroerr.ErrItemUpdateFailed
    }

    // Recalculate cart statistics
    err = l.svcCtx.CartStatsModel.RecalculateStats(l.ctx, uint64(in.UserId))
    if err != nil {
        return nil, zeroerr.ErrStatsUpdateFailed
    }

    return &cart.UpdateItemResponse{
        Success: true,
    }, nil
}
