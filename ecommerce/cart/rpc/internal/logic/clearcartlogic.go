package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
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

    return &cart.ClearCartResponse{
        Success: true,
    }, nil
}
