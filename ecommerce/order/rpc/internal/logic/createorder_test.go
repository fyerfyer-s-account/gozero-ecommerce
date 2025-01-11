package logic

import (
    "context"
    "flag"
    "testing"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
    "github.com/stretchr/testify/assert"
    "github.com/zeromicro/go-zero/core/conf"
)

func TestCreateOrderLogic_CreateOrder(t *testing.T) {
    configFile := flag.String("f", "../../etc/order.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    tests := []struct {
        name    string
        req     *order.CreateOrderRequest
        wantErr error
    }{
        {
            name: "create order successfully",
            req: &order.CreateOrderRequest{
                UserId:   1,
                Address:  "Test Address",
                Receiver: "Test User",
                Phone:    "1234567890",
                Items: []*order.OrderItemRequest{
                    {
                        ProductId: 1,
                        SkuId:    1,
                        Quantity: 2,
                    },
                },
            },
            wantErr: nil,
        },
        {
            name: "invalid user id",
            req: &order.CreateOrderRequest{
                UserId:   0,
                Address:  "Test Address",
                Receiver: "Test User",
                Phone:    "1234567890",
                Items: []*order.OrderItemRequest{
                    {
                        ProductId: 1,
                        SkuId:    1,
                        Quantity: 2,
                    },
                },
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "empty items",
            req: &order.CreateOrderRequest{
                UserId:   1,
                Address:  "Test Address",
                Receiver: "Test User",
                Phone:    "1234567890",
                Items:    []*order.OrderItemRequest{},
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewCreateOrderLogic(context.Background(), ctx)
            resp, err := l.CreateOrder(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.NotEmpty(t, resp.OrderNo)
                assert.Greater(t, resp.PayAmount, float64(0))
            }
        })
    }
}