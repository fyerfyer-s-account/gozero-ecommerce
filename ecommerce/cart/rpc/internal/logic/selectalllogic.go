package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
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

    return &cart.SelectAllResponse{
        Success: true,
    }, nil
}
