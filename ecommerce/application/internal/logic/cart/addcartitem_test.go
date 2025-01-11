package cart

import (
	"context"
	"errors"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestAddCartItemLogic_AddCartItem(t *testing.T) {
    mockCart := NewCart(t)
    svcCtx := &svc.ServiceContext{
        CartRpc: mockCart,
    }

    tests := []struct {
        name    string
        ctx     context.Context
        req     *types.CartItemReq
        mock    func()
        wantErr bool
    }{
        {
            name: "add item successfully",
            ctx: context.WithValue(context.Background(), "userId", int64(1)),
            req: &types.CartItemReq{
                ProductId: 1,
                SkuId:    1,
                Quantity: 2,
            },
            mock: func() {
                mockCart.EXPECT().AddItem(
                    mock.Anything,
                    &cart.AddItemRequest{
                        UserId:    1,
                        ProductId: 1,
                        SkuId:    1,
                        Quantity: 2,
                    },
                ).Return(&cart.AddItemResponse{Success: true}, nil)
            },
            wantErr: false,
        },
        {
            name: "invalid quantity",
            ctx: context.WithValue(context.Background(), "userId", int64(1)),
            req: &types.CartItemReq{
                ProductId: 1,
                SkuId:    1,
                Quantity: 0,
            },
            mock:    func() {},
            wantErr: true,
        },
        {
            name: "rpc error",
            ctx: context.WithValue(context.Background(), "userId", int64(1)),
            req: &types.CartItemReq{
                ProductId: 1,
                SkuId:    1,
                Quantity: 1,
            },
            mock: func() {
                mockCart.EXPECT().AddItem(
                    mock.Anything,
                    mock.Anything,
                ).Return(nil, errors.New("rpc error"))
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mock()

            l := NewAddCartItemLogic(tt.ctx, svcCtx)
            err := l.AddCartItem(tt.req)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}