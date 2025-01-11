package cart

import (
    "context"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
    cart "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

    "github.com/zeromicro/go-zero/core/logx"
)

type DeleteCartItemLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewDeleteCartItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCartItemLogic {
    return &DeleteCartItemLogic{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *DeleteCartItemLogic) DeleteCartItem(req *types.DeleteItemReq) error {
    userId := l.ctx.Value("userId").(int64)
    if userId <= 0 || req.Id <= 0 {
        return zeroerr.ErrInvalidParam
    }

    _, err := l.svcCtx.CartRpc.RemoveItem(l.ctx, &cart.RemoveItemRequest{
        UserId:    userId,
        ProductId: req.Id,
        SkuId:    req.SkuId,
    })

    return err
}