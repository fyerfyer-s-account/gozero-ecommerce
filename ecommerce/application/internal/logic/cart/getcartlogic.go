package cart

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	cart "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCartLogic {
	return &GetCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCartLogic) GetCart() (*types.CartInfo, error) {
    userId := l.ctx.Value("userId").(int64)

    cart, err := l.svcCtx.CartRpc.GetCart(l.ctx, &cart.GetCartRequest{
        UserId: userId,
    })
    if err != nil {
        if err == zeroerr.ErrCartNotFound {
            return nil, zeroerr.ErrCartNotFound
        }
        return nil, err
    }

    resp := &types.CartInfo{
        Items:         make([]types.CartItem, 0),
        TotalPrice:    0,
        TotalQuantity: 0,
        SelectedPrice: 0,
        SelectedCount: 0,
    }

    for _, item := range cart.Items {
        cartItem := types.CartItem{
            Id:          item.Id,
            ProductId:   item.ProductId,
            ProductName: item.ProductName,
            SkuId:       item.SkuId,
            SkuName:     item.SkuName,
            Image:       item.Image,
            Price:       item.Price,
            Quantity:    item.Quantity,
            Selected:    item.Selected,
            Stock:       item.Stock,
            CreatedAt:   item.CreatedAt,
        }

        resp.Items = append(resp.Items, cartItem)
        resp.TotalQuantity += item.Quantity
        resp.TotalPrice += item.Price * float64(item.Quantity)

        if item.Selected {
            resp.SelectedCount += item.Quantity
            resp.SelectedPrice += item.Price * float64(item.Quantity)
        }
    }

    return resp, nil
}

