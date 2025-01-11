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

func TestGetSelectedItemsLogic_GetSelectedItems(t *testing.T) {
    mockCart := NewCart(t)
    svcCtx := &svc.ServiceContext{
        CartRpc: mockCart,
    }

    tests := []struct {
        name     string
        ctx      context.Context
        mock     func()
        wantResp *types.SelectedItemsResp
        wantErr  error
    }{
        {
            name: "get selected items successfully",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            mock: func() {
                mockCart.EXPECT().GetSelectedItems(
                    mock.Anything,
                    &cart.GetSelectedItemsRequest{UserId: 1},
                ).Return(&cart.GetSelectedItemsResponse{
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
                    },
                    TotalPrice:    200.0,
                    TotalQuantity: 2,
                }, nil)

                mockCart.EXPECT().CheckStock(
                    mock.Anything,
                    &cart.CheckStockRequest{UserId: 1},
                ).Return(&cart.CheckStockResponse{
                    AllInStock: true,
                }, nil)
            },
            wantResp: &types.SelectedItemsResp{
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
                },
                TotalPrice:    200.0,
                TotalQuantity: 2,
                ValidStock:    true,
            },
            wantErr: nil,
        },
        {
            name: "no selected items",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            mock: func() {
                mockCart.EXPECT().GetSelectedItems(
                    mock.Anything,
                    &cart.GetSelectedItemsRequest{UserId: 1},
                ).Return(&cart.GetSelectedItemsResponse{
                    Items:         []*cart.CartItem{}, // 确保这里是空列表
                    TotalPrice:    0,
                    TotalQuantity: 0,
                }, nil)
            
                mockCart.EXPECT().CheckStock(
                    mock.Anything,
                    &cart.CheckStockRequest{UserId: 1},
                ).Return(&cart.CheckStockResponse{
                    AllInStock: true,
                }, nil)
            },
            wantResp: &types.SelectedItemsResp{
                Items:         []types.CartItem{},
                TotalPrice:    0,
                TotalQuantity: 0,
                ValidStock:    true,
            },
            wantErr: nil,
        },
        {
            name: "get selected items rpc error",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            mock: func() {
                mockCart.EXPECT().GetSelectedItems(
                    mock.Anything,
                    &cart.GetSelectedItemsRequest{UserId: 1},
                ).Return(nil, errors.New("rpc error"))
            },
            wantResp: nil,
            wantErr:  errors.New("rpc error"),
        },
        {
            name: "check stock rpc error",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            mock: func() {
                mockCart.EXPECT().GetSelectedItems(
                    mock.Anything,
                    &cart.GetSelectedItemsRequest{UserId: 1},
                ).Return(&cart.GetSelectedItemsResponse{
                    Items:         []*cart.CartItem{},
                    TotalPrice:    0,
                    TotalQuantity: 0,
                }, nil)
        
                mockCart.EXPECT().CheckStock(
                    mock.Anything,
                    &cart.CheckStockRequest{UserId: 1},
                ).Return(nil, errors.New("stock check error"))
            },
            wantResp: nil,
            wantErr:  errors.New("stock check error"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mock()

            l := NewGetSelectedItemsLogic(tt.ctx, svcCtx)
            resp, err := l.GetSelectedItems()

            if tt.wantErr != nil {
                assert.EqualError(t, err, tt.wantErr.Error())
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantResp, resp)
            }
        })
    }
}