package logic

import (
	"context"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveItemLogic {
	return &RemoveItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveItemLogic) RemoveItem(in *cart.RemoveItemRequest) (*cart.RemoveItemResponse, error) {
    if in.UserId <= 0 || in.ProductId <= 0 || in.SkuId <= 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    // Get cart item info for statistics update
    item, err := l.svcCtx.CartItemsModel.FindOneByUserIdSkuId(l.ctx, uint64(in.UserId), uint64(in.SkuId))
    if err != nil {
        return nil, zeroerr.ErrItemNotFound
    }

    // Delete cart item
    err = l.svcCtx.CartItemsModel.Delete(l.ctx, item.Id)
    if err != nil {
        return nil, zeroerr.ErrItemDeleteFailed
    }

    // Update cart statistics
    err = l.svcCtx.CartStatsModel.RecalculateStats(l.ctx, uint64(in.UserId))
    if err != nil {
        return nil, zeroerr.ErrStatsUpdateFailed
    }

    // After cart item deletion and statistics update
    items := []types.CartItem{
        {
            ProductID: in.ProductId,
            SkuID:     in.SkuId,
            Quantity:  0, // Set to 0 since item is removed
            Selected:  (item.Selected == 1),
            Price:     item.Price,
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
        // Don't return error as item is already removed
    }

    return &cart.RemoveItemResponse{
        Success: true,
    }, nil
}
