package cart

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelectCartItemsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSelectCartItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectCartItemsLogic {
	return &SelectCartItemsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelectCartItemsLogic) SelectCartItems(req *types.BatchOperateReq) error {
    userId := l.ctx.Value("userId").(int64)
    if userId <= 0 || len(req.ItemIds) == 0 {
        return zeroerr.ErrInvalidParam
    }

    cartInfo, err := l.svcCtx.CartRpc.GetCart(l.ctx, &cart.GetCartRequest{
        UserId: userId,
    })
    if err != nil {
        return err
    }

    // Create map for O(1) lookup
    itemMap := make(map[int64]*cart.CartItem)
    for _, item := range cartInfo.Items {
        itemMap[item.Id] = item
    }

    for _, itemId := range req.ItemIds {
        item, exists := itemMap[itemId]
        if !exists {
            continue
        }

        _, err := l.svcCtx.CartRpc.SelectItem(l.ctx, &cart.SelectItemRequest{
            UserId:    userId,
            ProductId: item.ProductId,
            SkuId:    item.SkuId,
        })
        if err != nil {
            return err
        }
    }

    return nil
}
