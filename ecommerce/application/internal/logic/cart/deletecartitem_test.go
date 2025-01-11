package cart

import (
    "context"
    "errors"
    "testing"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestDeleteCartItemLogic_DeleteCartItem(t *testing.T) {
    mockCart := NewCart(t)
    svcCtx := &svc.ServiceContext{
        CartRpc: mockCart,
    }

    tests := []struct {
        name    string
        ctx     context.Context
        req     *types.DeleteItemReq
        mock    func()
        wantErr bool
    }{
        {
            name: "delete item successfully",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            req: &types.DeleteItemReq{
                Id:    1,
                SkuId: 1,
            },
            mock: func() {
                mockCart.EXPECT().RemoveItem(
                    mock.Anything,
                    &cart.RemoveItemRequest{
                        UserId:    1,
                        ProductId: 1,
                        SkuId:    1,
                    },
                ).Return(&cart.RemoveItemResponse{Success: true}, nil)
            },
            wantErr: false,
        },
        {
            name: "invalid item id",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            req: &types.DeleteItemReq{
                Id:    0,
                SkuId: 1,
            },
            mock:    func() {},
            wantErr: true,
        },
        {
            name: "rpc error",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            req: &types.DeleteItemReq{
                Id:    1,
                SkuId: 1,
            },
            mock: func() {
                mockCart.EXPECT().RemoveItem(
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

            l := NewDeleteCartItemLogic(tt.ctx, svcCtx)
            err := l.DeleteCartItem(tt.req)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}