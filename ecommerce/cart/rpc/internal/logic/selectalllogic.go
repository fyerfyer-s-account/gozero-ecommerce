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

type SelectAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSelectAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectAllLogic {
	return &SelectAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SelectAllLogic) SelectAll(in *cart.SelectAllRequest) (*cart.SelectAllResponse, error) {
    if in.UserId <= 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    // Update all items to selected
    err := l.svcCtx.CartItemsModel.UpdateAllSelected(l.ctx, uint64(in.UserId), 1)
    if err != nil {
        return nil, zeroerr.ErrSelectAllFailed
    }

    // Recalculate cart statistics
    err = l.svcCtx.CartStatsModel.RecalculateStats(l.ctx, uint64(in.UserId))
    if err != nil {
        return nil, zeroerr.ErrStatsUpdateFailed
    }

    // Get all cart items after update
    items, err := l.svcCtx.CartItemsModel.FindByUserId(l.ctx, uint64(in.UserId))
    if err != nil {
        return nil, err
    }

    // Convert cart items to event items
    eventItems := make([]types.CartItem, 0, len(items))
    for _, item := range items {
        eventItems = append(eventItems, types.CartItem{
            ProductID: int64(item.ProductId),
            SkuID:     int64(item.SkuId),
            Quantity:  int32(item.Quantity),
            Selected:  true,
            Price:     item.Price,
        })
    }

    selectionEvent := &types.CartSelectionEvent{
        CartEvent: types.CartEvent{
            Type:      types.CartSelected,
            UserID:    in.UserId,
            Timestamp: time.Now(),
        },
        Items: eventItems,
    }

    if err := l.svcCtx.Producer.PublishCartSelected(l.ctx, selectionEvent); err != nil {
        logx.Errorf("Failed to publish cart selection event: %v", err)
        // Don't return error as items are already selected
    }

    return &cart.SelectAllResponse{
        Success: true,
    }, nil
}
