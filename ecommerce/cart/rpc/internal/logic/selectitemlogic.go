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

type SelectItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSelectItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectItemLogic {
	return &SelectItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 商品选择
func (l *SelectItemLogic) SelectItem(in *cart.SelectItemRequest) (*cart.SelectItemResponse, error) {
    if in.UserId <= 0 || in.ProductId <= 0 || in.SkuId <= 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    // Update selected status
    err := l.svcCtx.CartItemsModel.UpdateSelected(l.ctx, uint64(in.UserId), uint64(in.ProductId), uint64(in.SkuId), 1)
    if err != nil {
        return nil, zeroerr.ErrSelectFailed
    }

    // Recalculate cart stats
    err = l.svcCtx.CartStatsModel.RecalculateStats(l.ctx, uint64(in.UserId))
    if err != nil {
        return nil, zeroerr.ErrStatsUpdateFailed
    }

    // Get item details after update
    item, err := l.svcCtx.CartItemsModel.FindOneByUserIdSkuId(l.ctx, uint64(in.UserId), uint64(in.SkuId))
    if err != nil {
        return nil, zeroerr.ErrItemNotFound
    }

    // Create and publish cart selection event
    items := []types.CartItem{
        {
            ProductID: in.ProductId,
            SkuID:     in.SkuId,
            Quantity:  int32(item.Quantity),
            Selected:  true,
            Price:     item.Price,
        },
    }

    selectionEvent := &types.CartSelectionEvent{
        CartEvent: types.CartEvent{
            Type:      types.CartSelected,
            UserID:    in.UserId,
            Timestamp: time.Now(),
        },
        Items: items,
    }

    if err := l.svcCtx.Producer.PublishCartSelected(l.ctx, selectionEvent); err != nil {
        logx.Errorf("Failed to publish cart selection event: %v", err)
        // Don't return error as cart item is already selected
    }

    return &cart.SelectItemResponse{
        Success: true,
    }, nil
}
