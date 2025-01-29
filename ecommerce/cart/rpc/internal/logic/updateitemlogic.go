package logic

import (
	"context"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
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
	item, err := l.svcCtx.CartItemsModel.FindOneByUserIdSkuId(l.ctx, uint64(in.UserId), uint64(in.SkuId))
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

	// After cart item and statistics are updated successfully
    items := []types.CartItem{
        {
            ProductID: in.ProductId,
            SkuID:     in.SkuId,
            Quantity:  in.Quantity,
            Selected:  (item.Selected == 1),
            Price:     sku.Sku.Price,
        },
    }

    updatedEvent := &types.CartUpdatedEvent{
        CartEvent: types.CartEvent{
            Type:      types.CartUpdated,
            UserID:    in.UserId,
            Timestamp: time.Now(),
        },
        Items: items,
    }

    if err := l.svcCtx.Producer.PublishCartUpdated(l.ctx, updatedEvent); err != nil {
        logx.Errorf("Failed to publish cart updated event: %v", err)
        // Don't return error as cart item is already updated
    }

	return &cart.UpdateItemResponse{
		Success: true,
	}, nil
}
