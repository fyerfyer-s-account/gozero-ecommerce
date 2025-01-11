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

func TestGetCartLogic_GetCart(t *testing.T) {
    mockCart := NewCart(t)
    svcCtx := &svc.ServiceContext{
        CartRpc: mockCart,
    }

    tests := []struct {
        name     string
        ctx      context.Context
        mock     func()
        wantResp *types.CartInfo
        wantErr  error
    }{
        {
            name: "get cart successfully with items",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            mock: func() {
                mockCart.EXPECT().GetCart(
                    mock.Anything,
                    &cart.GetCartRequest{UserId: 1},
                ).Return(&cart.GetCartResponse{
                    Items: []*cart.CartItem{
                        {
                            Id:          1,
                            ProductId:   100,
                            ProductName: "Product 1",
                            SkuId:       200,
                            SkuName:     "SKU 1",
                            Image:       "image1.jpg",
                            Price:       100.0,
                            Quantity:    2,
                            Selected:    true,
                            Stock:       10,
                            CreatedAt:   1234567890,
                        },
                        {
                            Id:          2,
                            ProductId:   101,
                            ProductName: "Product 2",
                            SkuId:       201,
                            SkuName:     "SKU 2",
                            Image:       "image2.jpg",
                            Price:       50.0,
                            Quantity:    1,
                            Selected:    false,
                            Stock:       5,
                            CreatedAt:   1234567891,
                        },
                    },
                }, nil)
            },
            wantResp: &types.CartInfo{
                Items: []types.CartItem{
                    {
                        Id:          1,
                        ProductId:   100,
                        ProductName: "Product 1",
                        SkuId:       200,
                        SkuName:     "SKU 1",
                        Image:       "image1.jpg",
                        Price:       100.0,
                        Quantity:    2,
                        Selected:    true,
                        Stock:       10,
                        CreatedAt:   1234567890,
                    },
                    {
                        Id:          2,
                        ProductId:   101,
                        ProductName: "Product 2",
                        SkuId:       201,
                        SkuName:     "SKU 2",
                        Image:       "image2.jpg",
                        Price:       50.0,
                        Quantity:    1,
                        Selected:    false,
                        Stock:       5,
                        CreatedAt:   1234567891,
                    },
                },
                TotalPrice:    250.0,  // (100*2 + 50*1)
                TotalQuantity: 3,      // (2 + 1)
                SelectedPrice: 200.0,  // (100*2)
                SelectedCount: 2,      // (2)
            },
            wantErr: nil,
        },
        {
            name: "get empty cart",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            mock: func() {
                mockCart.EXPECT().GetCart(
                    mock.Anything,
                    &cart.GetCartRequest{UserId: 1},
                ).Return(&cart.GetCartResponse{
                    Items: []*cart.CartItem{},
                }, nil)
            },
            wantResp: &types.CartInfo{
                Items:         []types.CartItem{},
                TotalPrice:   0,
                TotalQuantity: 0,
                SelectedPrice: 0,
                SelectedCount: 0,
            },
            wantErr: nil,
        },
        {
            name: "get cart rpc error",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            mock: func() {
                mockCart.EXPECT().GetCart(
                    mock.Anything,
                    &cart.GetCartRequest{UserId: 1},
                ).Return(nil, errors.New("rpc error"))
            },
            wantResp: nil,
            wantErr:  errors.New("rpc error"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mock()

            l := NewGetCartLogic(tt.ctx, svcCtx)
            resp, err := l.GetCart()

            if tt.wantErr != nil {
                assert.EqualError(t, err, tt.wantErr.Error())
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantResp, resp)
            }

            mockCart.AssertExpectations(t)
        })
    }
}