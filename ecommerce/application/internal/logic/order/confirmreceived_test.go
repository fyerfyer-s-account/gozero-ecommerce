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

func TestConfirmReceivedLogic_ConfirmReceived(t *testing.T) {
    mockOrder := NewOrderService(t)
    svcCtx := &svc.ServiceContext{
        OrderRpc: mockOrder,
    }

    tests := []struct {
        name    string
        req     *types.ConfirmOrderReq
        mock    func()
        wantErr error
    }{
        {
            name: "confirm receive successfully",
            req: &types.ConfirmOrderReq{
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
                        Status:  3, // Shipped status
                    },
                }, nil)

                mockOrder.On("ReceiveOrder",
                    mock.Anything,
                    &order.ReceiveOrderRequest{
                        OrderNo: "ORDER123",
                    },
                ).Return(&order.ReceiveOrderResponse{
                    Success: true,
                }, nil)
            },
            wantErr: nil,
        },
        {
            name: "order not found",
            req: &types.ConfirmOrderReq{
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

            l := NewConfirmReceivedLogic(context.Background(), svcCtx)
            err := l.ConfirmReceived(tt.req)

            if tt.wantErr != nil {
                assert.Equal(t, tt.wantErr, err)
            } else {
                assert.Nil(t, err)
            }
        })
    }
}