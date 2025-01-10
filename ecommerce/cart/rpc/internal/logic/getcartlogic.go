package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCartLogic {
	return &GetCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCartLogic) GetCart(in *cart.GetCartRequest) (*cart.GetCartResponse, error) {
    if in.UserId <= 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    items, err := l.svcCtx.CartItemsModel.FindByUserId(l.ctx, uint64(in.UserId))
    if err != nil {
        return nil, err
    }

    stats, err := l.svcCtx.CartStatsModel.FindOne(l.ctx, uint64(in.UserId))
    if err != nil && err != model.ErrNotFound {
        return nil, err
    }

    cartItems := make([]*cart.CartItem, 0, len(items))
    for _, item := range items {
        cartItems = append(cartItems, &cart.CartItem{
            Id:          int64(item.Id),
            UserId:      int64(item.UserId),
            ProductId:   int64(item.ProductId),
            SkuId:       int64(item.SkuId),
            ProductName: item.ProductName,
            SkuName:     item.SkuName,
            Image:       item.Image.String,
            Price:       item.Price,
            Quantity:    item.Quantity,
            Selected:    item.Selected == 1,
            CreatedAt:   item.CreatedAt.Unix(),
            UpdatedAt:   item.UpdatedAt.Unix(),
        })
    }

    resp := &cart.GetCartResponse{
        Items:         cartItems,
        TotalQuantity: 0,
        TotalPrice:    0,
    }

    if stats != nil {
        resp.TotalQuantity = int32(stats.TotalQuantity)
        resp.TotalPrice = stats.TotalAmount
    }

    return resp, nil
}
