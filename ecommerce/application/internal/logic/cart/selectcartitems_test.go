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

func TestSelectCartItemsLogic_SelectCartItems(t *testing.T) {
    mockCart := NewCart(t)
    svcCtx := &svc.ServiceContext{
        CartRpc: mockCart,
    }

    tests := []struct {
        name    string
        ctx     context.Context
        req     *types.BatchOperateReq
        mock    func()
        wantErr error
    }{
        {
            name: "select items successfully",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            req:  &types.BatchOperateReq{ItemIds: []int64{1, 2}},
            mock: func() {
                mockCart.EXPECT().GetCart(
                    mock.Anything,
                    &cart.GetCartRequest{UserId: 1},
                ).Return(&cart.GetCartResponse{
                    Items: []*cart.CartItem{
                        {Id: 1, ProductId: 100, SkuId: 200},
                        {Id: 2, ProductId: 101, SkuId: 201},
                    },
                }, nil)

                mockCart.EXPECT().SelectItem(
                    mock.Anything,
                    &cart.SelectItemRequest{UserId: 1, ProductId: 100, SkuId: 200},
                ).Return(&cart.SelectItemResponse{Success: true}, nil)

                mockCart.EXPECT().SelectItem(
                    mock.Anything,
                    &cart.SelectItemRequest{UserId: 1, ProductId: 101, SkuId: 201},
                ).Return(&cart.SelectItemResponse{Success: true}, nil)
            },
            wantErr: nil,
        },
        {
            name: "empty item ids",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            req:  &types.BatchOperateReq{ItemIds: []int64{}},
            mock: func() {},
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "get cart rpc error",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            req:  &types.BatchOperateReq{ItemIds: []int64{1}},
            mock: func() {
                mockCart.EXPECT().GetCart(
                    mock.Anything,
                    &cart.GetCartRequest{UserId: 1},
                ).Return(nil, errors.New("rpc error"))
            },
            wantErr: errors.New("rpc error"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mock()

            l := NewSelectCartItemsLogic(tt.ctx, svcCtx)
            err := l.SelectCartItems(tt.req)

            if tt.wantErr == nil {
                assert.NoError(t, err)
            } else {
                assert.EqualError(t, err, tt.wantErr.Error())
            }

            mockCart.AssertExpectations(t)
        })
    }
}