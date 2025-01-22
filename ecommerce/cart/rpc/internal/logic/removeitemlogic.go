package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/google/uuid"

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

    event := &types.CartEvent{
        ID:        uuid.New().String(),
        Type:      types.EventTypeItemRemoved,
        Timestamp: time.Now(),
        Data: &types.CartItemRemovedData{
            UserID:    in.UserId,
            ProductID: in.ProductId,
        },
        Metadata: types.EventMetadata{
            TraceID: l.ctx.Value("trace_id").(string),
            UserID:  fmt.Sprint(in.UserId),
        },
    }

    if err := l.svcCtx.Producer.PublishEvent(event); err != nil {
        logx.Errorf("Failed to publish cart.item.removed event: %v", err)
    }

    return &cart.RemoveItemResponse{
        Success: true,
    }, nil
}
