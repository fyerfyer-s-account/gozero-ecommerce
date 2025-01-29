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

type UnselectAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnselectAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnselectAllLogic {
	return &UnselectAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnselectAllLogic) UnselectAll(in *cart.UnselectAllRequest) (*cart.UnselectAllResponse, error) {
    if in.UserId <= 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    err := l.svcCtx.CartItemsModel.UpdateAllSelected(l.ctx, uint64(in.UserId), 0)
    if err != nil {
        return nil, zeroerr.ErrDeselectFailed
    }

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
            Selected:  false,
            Price:     item.Price,
        })
    }

    selectionEvent := &types.CartSelectionEvent{
        CartEvent: types.CartEvent{
            Type:      types.CartUnselected,
            UserID:    in.UserId,
            Timestamp: time.Now(),
        },
        Items: eventItems,
    }

    if err := l.svcCtx.Producer.PublishCartSelected(l.ctx, selectionEvent); err != nil {
        logx.Errorf("Failed to publish cart unselection event: %v", err)
        // Don't return error as items are already unselected
    }

    return &cart.UnselectAllResponse{
        Success: true,
    }, nil
}
