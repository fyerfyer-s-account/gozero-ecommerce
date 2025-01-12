package order

import (
    "context"
    "testing"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/orderservice"
    "github.com/stretchr/testify/assert"
    mock "github.com/stretchr/testify/mock"
)

func TestListOrdersLogic_ListOrders(t *testing.T) {
    mockOrder := NewOrderService(t)
    svcCtx := &svc.ServiceContext{
        OrderRpc: mockOrder,
    }

    tests := []struct {
        name    string
        ctx     context.Context
        req     *types.OrderListReq
        mock    func()
        want    *types.OrderListResp
        wantErr error
    }{
        {
            name: "list orders successfully",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            req: &types.OrderListReq{
                Status:   1,
                Page:     1,
                PageSize: 10,
            },
            mock: func() {
                mockOrder.On("ListOrders",
                    mock.Anything,
                    &orderservice.ListOrdersRequest{
                        UserId:   1,
                        Status:   1,
                        Page:     1,
                        PageSize: 10,
                    },
                ).Return(&orderservice.ListOrdersResponse{
                    Orders: []*orderservice.Order{
                        {
                            Id:          1,
                            OrderNo:     "ORDER123",
                            UserId:      1,
                            Status:      1,
                            TotalAmount: 100,
                            Items: []*orderservice.OrderItem{
                                {
                                    ProductId:   1,
                                    ProductName: "Test Product",
                                    Quantity:    1,
                                    Price:      100,
                                },
                            },
                        },
                    },
                    Total: 1,
                }, nil)
            },
            want: &types.OrderListResp{
                List: []types.Order{
                    {
                        Id:          1,
                        OrderNo:     "ORDER123",
                        UserId:      1,
                        Status:      1,
                        TotalAmount: 100,
                        Items: []types.OrderItem{
                            {
                                ProductId:   1,
                                ProductName: "Test Product",
                                Quantity:    1,
                                Price:      100,
                            },
                        },
                    },
                },
                Total:      1,
                Page:       1,
                TotalPages: 1,
            },
            wantErr: nil,
        },
        {
            name: "empty result",
            ctx:  context.WithValue(context.Background(), "userId", int64(1)),
            req: &types.OrderListReq{
                Status:   1,
                Page:     1,
                PageSize: 10,
            },
            mock: func() {
                mockOrder.On("ListOrders",
                    mock.Anything,
                    &orderservice.ListOrdersRequest{
                        UserId:   1,
                        Status:   1,
                        Page:     1,
                        PageSize: 10,
                    },
                ).Return(&orderservice.ListOrdersResponse{
                    Orders: []*orderservice.Order{},
                    Total:  0,
                }, nil)
            },
            want: &types.OrderListResp{
                List:       []types.Order{},
                Total:      0,
                Page:       1,
                TotalPages: 0,
            },
            wantErr: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mock()

            l := NewListOrdersLogic(tt.ctx, svcCtx)
            got, err := l.ListOrders(tt.req)

            if tt.wantErr != nil {
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, got)
            } else {
                assert.Nil(t, err)
                assert.Equal(t, tt.want.Total, got.Total)
                assert.Equal(t, tt.want.Page, got.Page)
                assert.Equal(t, tt.want.TotalPages, got.TotalPages)
                assert.Equal(t, len(tt.want.List), len(got.List))
            }
        })
    }
}