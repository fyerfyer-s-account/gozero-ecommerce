package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSelectedItemsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSelectedItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSelectedItemsLogic {
	return &GetSelectedItemsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 结算相关
func (l *GetSelectedItemsLogic) GetSelectedItems(in *cart.GetSelectedItemsRequest) (*cart.GetSelectedItemsResponse, error) {
    if in.UserId <= 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    // Get selected items
    items, err := l.svcCtx.CartItemsModel.FindSelectedByUserId(l.ctx, uint64(in.UserId))
    if err != nil {
        return nil, err
    }

    // Get statistics
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
            SkuName:    item.SkuName,
            Image:      item.Image.String,
            Price:      item.Price,
            Quantity:   item.Quantity,
            Selected:   item.Selected == 1,
            CreatedAt:  item.CreatedAt.Unix(),
            UpdatedAt:  item.UpdatedAt.Unix(),
        })
    }

    resp := &cart.GetSelectedItemsResponse{
        Items:         cartItems,
        TotalQuantity: 0,
        TotalPrice:    0,
    }

    if stats != nil {
        resp.TotalQuantity = int32(stats.SelectedQuantity)
        resp.TotalPrice = stats.SelectedAmount
    }

    return resp, nil
}
