package cart

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCartItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCartItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCartItemLogic {
	return &UpdateCartItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCartItemLogic) UpdateCartItem(req *types.CartItemReq) error {
    userId := l.ctx.Value("userId").(int64)
    if userId <= 0 || req.ProductId <= 0 || req.SkuId <= 0 || req.Quantity <= 0 {
        return zeroerr.ErrInvalidParam
    }

    _, err := l.svcCtx.CartRpc.UpdateItem(l.ctx, &cart.UpdateItemRequest{
        UserId:    userId,
        ProductId: req.ProductId,
        SkuId:    req.SkuId,
        Quantity: int32(req.Quantity),
    })

    return err
}
