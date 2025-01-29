package logic

import (
	"context"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearCartLogic {
	return &ClearCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ClearCartLogic) ClearCart(in *cart.ClearCartRequest) (*cart.ClearCartResponse, error) {
    if in.UserId <= 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    // Delete all cart items
    err := l.svcCtx.CartItemsModel.DeleteByUserId(l.ctx, uint64(in.UserId))
    if err != nil {
        return nil, zeroerr.ErrCartDeleteFailed
    }

    // Reset cart statistics
    stats := &model.CartStatistics{
        UserId:           uint64(in.UserId),
        TotalQuantity:    0,
        SelectedQuantity: 0,
        TotalAmount:      0,
        SelectedAmount:   0,
    }
    err = l.svcCtx.CartStatsModel.Upsert(l.ctx, stats)
    if err != nil {
        return nil, zeroerr.ErrCartUpdateFailed
    }

    // After cart items deletion and statistics update
    clearedEvent := &types.CartClearedEvent{
        CartEvent: types.CartEvent{
            Type:      types.CartCleared,
            UserID:    in.UserId,
            Timestamp: time.Now(),
        },
        Reason: "user_requested",
    }

    if err := l.svcCtx.Producer.PublishCartCleared(l.ctx, clearedEvent); err != nil {
        logx.Errorf("Failed to publish cart cleared event: %v", err)
        // Don't return error as cart is already cleared
    }

    return &cart.ClearCartResponse{
        Success: true,
    }, nil
}