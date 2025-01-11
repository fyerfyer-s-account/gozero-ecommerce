package cart

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	cart "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSelectedItemsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSelectedItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSelectedItemsLogic {
	return &GetSelectedItemsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSelectedItemsLogic) GetSelectedItems() (resp *types.SelectedItemsResp, err error) {
    userId := l.ctx.Value("userId").(int64)

    selectedItems, err := l.svcCtx.CartRpc.GetSelectedItems(l.ctx, &cart.GetSelectedItemsRequest{
        UserId: userId,
    })
    if err != nil {
        return nil, err
    }

    resp = &types.SelectedItemsResp{
        Items:         make([]types.CartItem, 0),
        TotalPrice:    selectedItems.TotalPrice,
        TotalQuantity: int64(selectedItems.TotalQuantity),
        ValidStock:    true,
    }

    // Move stock check after error handling
    stockCheck, err := l.svcCtx.CartRpc.CheckStock(l.ctx, &cart.CheckStockRequest{
        UserId: userId,
    })
    if err != nil {
        return nil, err
    }

    resp.ValidStock = stockCheck.AllInStock

    // Only process selected items
    for _, item := range selectedItems.Items {
        if !item.Selected {
            continue
        }
        cartItem := types.CartItem{
            Id:          item.Id,
            ProductId:   item.ProductId,
            ProductName: item.ProductName,
            SkuId:      item.SkuId,
            SkuName:    item.SkuName,
            Image:      item.Image,
            Price:      item.Price,
            Quantity:   item.Quantity,
            Selected:   item.Selected,
            Stock:      item.Stock,
            CreatedAt:  item.CreatedAt,
        }
        resp.Items = append(resp.Items, cartItem)
    }

    return resp, nil
}