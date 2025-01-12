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

func TestCancelOrderLogic_CancelOrder(t *testing.T) {
    mockOrder := NewOrderService(t)
    svcCtx := &svc.ServiceContext{
        OrderRpc: mockOrder,
    }

    tests := []struct {
        name    string
        ctx     context.Context
        req     *types.CancelOrderReq
        mock    func()
        wantErr error
    }{
        {
            name: "cancel order successfully",
            ctx:  context.Background(),
            req: &types.CancelOrderReq{
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
                        OrderNo: "ORDER123",
                    },
                }, nil)

                mockOrder.On("CancelOrder",
                    mock.Anything,
                    &order.CancelOrderRequest{
                        OrderNo: "ORDER123",
                        Reason:  "用户取消",
                    },
                ).Return(&order.CancelOrderResponse{
                    Success: true,
                }, nil)
            },
            wantErr: nil,
        },
        {
            name: "order not found",
            ctx:  context.Background(),
            req: &types.CancelOrderReq{
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
            wantErr: zeroerr.ErrOrderNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mock()

            l := NewCancelOrderLogic(tt.ctx, svcCtx)
            err := l.CancelOrder(tt.req)

            if tt.wantErr != nil {
                assert.Equal(t, tt.wantErr, err)
            } else {
                assert.Nil(t, err)
            }
        })
    }
}