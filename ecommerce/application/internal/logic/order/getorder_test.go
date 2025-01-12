package order

import (
    "context"
    "testing"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
    "github.com/stretchr/testify/assert"
    mock "github.com/stretchr/testify/mock"
)

func TestGetOrderLogic_GetOrder(t *testing.T) {
    mockOrder := NewOrderService(t)
    svcCtx := &svc.ServiceContext{
        OrderRpc: mockOrder,
    }

    tests := []struct {
        name    string
        req     *types.GetOrderReq
        mock    func()
        want    *types.Order
        wantErr error
    }{
        {
            name: "get order successfully",
            req: &types.GetOrderReq{
                Id: 123,
            },
            mock: func() {
                mockOrder.On("GetOrder",
                    mock.Anything,
                    &order.GetOrderRequest{
                        OrderNo: "123",
                    },
                ).Return(&order.GetOrderResponse{
                    Order: &order.Order{
                        Id:          123,
                        OrderNo:     "ORDER123",
                        UserId:      1,
                        Status:      1,
                        TotalAmount: 100,
                        Items: []*order.OrderItem{{
                            Id:          1,
                            ProductId:   1,
                            ProductName: "Test Product",
                            SkuId:      1,
                            SkuName:    "Test SKU",
                            Price:      50,
                            Quantity:   2,
                            TotalAmount: 100,
                        }},
                    },
                }, nil)
            },
            want: &types.Order{
                Id:          123,
                OrderNo:     "ORDER123",
                UserId:      1,
                Status:      1,
                TotalAmount: 100,
                Items: []types.OrderItem{{
                    Id:          1,
                    ProductId:   1,
                    ProductName: "Test Product",
                    SkuId:      1,
                    SkuName:    "Test SKU",
                    Price:      50,
                    Quantity:   2,
                    Amount:     100,
                }},
            },
            wantErr: nil,
        },
        {
            name: "order not found",
            req: &types.GetOrderReq{
                Id: 999,
            },
            mock: func() {
                mockOrder.On("GetOrder",
                    mock.Anything,
                    &order.GetOrderRequest{
                        OrderNo: "999",
                    },
                ).Return(nil, zeroerr.ErrOrderNotFound)
            },
            want:    nil,
            wantErr: zeroerr.ErrOrderNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mock()

            l := NewGetOrderLogic(context.Background(), svcCtx)
            got, err := l.GetOrder(tt.req)

            if tt.wantErr != nil {
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, got)
            } else {
                assert.Nil(t, err)
                assert.Equal(t, tt.want.Id, got.Id)
                assert.Equal(t, tt.want.OrderNo, got.OrderNo)
                assert.Equal(t, tt.want.Status, got.Status)
                assert.Equal(t, tt.want.TotalAmount, got.TotalAmount)
                assert.Equal(t, len(tt.want.Items), len(got.Items))
            }
        })
    }
}