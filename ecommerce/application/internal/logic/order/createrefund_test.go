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

func TestCreateRefundLogic_CreateRefund(t *testing.T) {
    mockOrder := NewOrderService(t)
    svcCtx := &svc.ServiceContext{
        OrderRpc: mockOrder,
    }

    tests := []struct {
        name    string
        req     *types.CreateRefundReq
        mock    func()
        want    *types.RefundInfo
        wantErr error
    }{
        {
            name: "create refund successfully",
            req: &types.CreateRefundReq{
                Id:      1,
                OrderNo: "ORDER123",
                Amount:  50,
                Reason:  "test refund",
                Desc:    "test description",
                Images:  []string{"image1.jpg"},
            },
            mock: func() {
                mockOrder.On("GetOrder",
                    mock.Anything,
                    &order.GetOrderRequest{
                        OrderNo: "ORDER123",
                    },
                ).Return(&order.GetOrderResponse{
                    Order: &order.Order{
                        OrderNo:   "ORDER123",
                        PayAmount: 100,
                    },
                }, nil)

                mockOrder.On("CreateRefund",
                    mock.Anything,
                    &order.CreateRefundRequest{
                        OrderNo:     "ORDER123",
                        Amount:      50,
                        Reason:      "test refund",
                        Description: "test description",
                        Images:      []string{"image1.jpg"},
                    },
                ).Return(&order.CreateRefundResponse{
                    RefundNo: "RF123",
                }, nil)
            },
            want: &types.RefundInfo{
                Id:       1,
                RefundNo: "RF123",
                Status:   0,
                Amount:   50,
                Reason:  "test refund",
                Desc:    "test description",
                Images:  []string{"image1.jpg"},
            },
            wantErr: nil,
        },
        {
            name: "refund amount exceeds order amount",
            req: &types.CreateRefundReq{
                OrderNo: "ORDER123",
                Amount:  150,
            },
            mock: func() {
                mockOrder.On("GetOrder",
                    mock.Anything,
                    &order.GetOrderRequest{
                        OrderNo: "ORDER123",
                    },
                ).Return(&order.GetOrderResponse{
                    Order: &order.Order{
                        OrderNo:   "ORDER123",
                        PayAmount: 100,
                    },
                }, nil)
            },
            want:    nil,
            wantErr: zeroerr.ErrRefundExceedAmount,
        },
        {
            name: "order not found",
            req: &types.CreateRefundReq{
                OrderNo: "NOT_EXIST",
                Amount:  50,
            },
            mock: func() {
                mockOrder.On("GetOrder",
                    mock.Anything,
                    &order.GetOrderRequest{
                        OrderNo: "NOT_EXIST",
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

            l := NewCreateRefundLogic(context.Background(), svcCtx)
            got, err := l.CreateRefund(tt.req)

            if tt.wantErr != nil {
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, got)
            } else {
                assert.Nil(t, err)
                assert.Equal(t, tt.want.Id, got.Id)
                assert.Equal(t, tt.want.RefundNo, got.RefundNo)
                assert.Equal(t, tt.want.Amount, got.Amount)
                assert.Equal(t, tt.want.Status, got.Status)
                assert.Equal(t, tt.want.Reason, got.Reason)
                assert.Equal(t, tt.want.Desc, got.Desc)
                assert.Equal(t, tt.want.Images, got.Images)
            }
        })
    }
}