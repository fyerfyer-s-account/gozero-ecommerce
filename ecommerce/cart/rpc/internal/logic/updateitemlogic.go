package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/google/uuid"

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

	event := &types.CartEvent{
		ID:        uuid.New().String(),
		Type:      types.EventTypeItemUpdated,
		Timestamp: time.Now(),
		Data: &types.CartItemUpdatedData{
			CartItemData: types.CartItemData{
				UserID:    in.UserId,
				ProductID: in.ProductId,
				Quantity:  in.Quantity,
			},
			OldQuantity: int32(item.Quantity),
		},
		Metadata: types.EventMetadata{
			TraceID: l.ctx.Value("trace_id").(string),
			UserID:  fmt.Sprint(in.UserId),
		},
	}

	if err := l.svcCtx.Producer.PublishEvent(event); err != nil {
		logx.Errorf("Failed to publish cart.item.updated event: %v", err)
	}

	return &cart.UpdateItemResponse{
		Success: true,
	}, nil
}
