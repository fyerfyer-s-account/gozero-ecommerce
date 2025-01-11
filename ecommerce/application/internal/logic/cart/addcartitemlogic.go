package cart

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	cart "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCartItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCartItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCartItemLogic {
	return &AddCartItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCartItemLogic) AddCartItem(req *types.CartItemReq) error {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("userId").(int64)
    if userId <= 0 || req.ProductId <= 0 || req.SkuId <= 0 || req.Quantity <= 0 {
        return zeroerr.ErrInvalidParam
    }

    _, err := l.svcCtx.CartRpc.AddItem(l.ctx, &cart.AddItemRequest{
        UserId:    userId,
        ProductId: req.ProductId,
        SkuId:    req.SkuId,
        Quantity: int64(req.Quantity),
    })

    return err
}
