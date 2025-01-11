package cart

import (
    "context"
    "errors"
    "testing"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestUpdateCartItemLogic_UpdateCartItem(t *testing.T) {
    mockCart := NewCart(t)
    svcCtx := &svc.ServiceContext{
        CartRpc: mockCart,
    }

    tests := []struct {
        name    string
        ctx     context.Context
        req     *types.CartItemReq
        mock    func()
        wantErr error
    }{
        {
            name: "update cart item successfully",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            req: &types.CartItemReq{
                ProductId: 100,
                SkuId:    200,
                Quantity: 2,
            },
            mock: func() {
                mockCart.EXPECT().UpdateItem(
                    mock.Anything,
                    &cart.UpdateItemRequest{
                        UserId:    1,
                        ProductId: 100,
                        SkuId:    200,
                        Quantity: 2,
                    },
                ).Return(&cart.UpdateItemResponse{Success: true}, nil)
            },
            wantErr: nil,
        },
        {
            name: "invalid parameters",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            req: &types.CartItemReq{
                ProductId: 0,
                SkuId:    200,
                Quantity: 2,
            },
            mock:    func() {},
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "rpc error",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            req: &types.CartItemReq{
                ProductId: 100,
                SkuId:    200,
                Quantity: 2,
            },
            mock: func() {
                mockCart.EXPECT().UpdateItem(
                    mock.Anything,
                    &cart.UpdateItemRequest{
                        UserId:    1,
                        ProductId: 100,
                        SkuId:    200,
                        Quantity: 2,
                    },
                ).Return(nil, errors.New("rpc error"))
            },
            wantErr: errors.New("rpc error"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mock()

            l := NewUpdateCartItemLogic(tt.ctx, svcCtx)
            err := l.UpdateCartItem(tt.req)

            if tt.wantErr == nil {
                assert.NoError(t, err)
            } else {
                assert.EqualError(t, err, tt.wantErr.Error())
            }

            mockCart.AssertExpectations(t)
        })
    }
}