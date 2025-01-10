package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnselectItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnselectItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnselectItemLogic {
	return &UnselectItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnselectItemLogic) UnselectItem(in *cart.UnselectItemRequest) (*cart.UnselectItemResponse, error) {
    if in.UserId <= 0 || in.ProductId <= 0 || in.SkuId <= 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    err := l.svcCtx.CartItemsModel.UpdateSelected(l.ctx, uint64(in.UserId), uint64(in.ProductId), uint64(in.SkuId), 0)
    if err != nil {
        return nil, zeroerr.ErrDeselectFailed
    }

    err = l.svcCtx.CartStatsModel.RecalculateStats(l.ctx, uint64(in.UserId))
    if err != nil {
        return nil, zeroerr.ErrStatsUpdateFailed
    }

    return &cart.UnselectItemResponse{
        Success: true,
    }, nil
}
