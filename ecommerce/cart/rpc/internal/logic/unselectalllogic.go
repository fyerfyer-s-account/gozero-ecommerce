package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
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

    return &cart.UnselectAllResponse{
        Success: true,
    }, nil
}
