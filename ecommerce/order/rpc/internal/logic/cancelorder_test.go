package logic

import (
    "context"
    "flag"
    "testing"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
    "github.com/stretchr/testify/assert"
    "github.com/zeromicro/go-zero/core/conf"
)

func TestCancelOrderLogic_CancelOrder(t *testing.T) {
    configFile := flag.String("f", "../../etc/order.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test order
    testOrder := &model.Orders{
        OrderNo:        "TEST_ORDER_001",
        UserId:        1,
        TotalAmount:   100,
        PayAmount:     100,
        FreightAmount: 0,
        Status:        1, // Pending payment
        Address:       "Test Address",
        Receiver:      "Test User",
        Phone:         "1234567890",
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }
    
    result, err := ctx.OrdersModel.Insert(context.Background(), testOrder)
    assert.NoError(t, err)
    orderId, err := result.LastInsertId()
    assert.NoError(t, err)

    tests := []struct {
        name    string
        req     *order.CancelOrderRequest
        wantErr error
    }{
        {
            name: "cancel order successfully",
            req: &order.CancelOrderRequest{
                OrderNo: "TEST_ORDER_001",
                Reason:  "test cancel",
            },
            wantErr: nil,
        },
        {
            name: "invalid order number",
            req: &order.CancelOrderRequest{
                OrderNo: "",
                Reason:  "test cancel",
            },
            wantErr: zeroerr.ErrOrderInvalidParam,
        },
        {
            name: "order not found",
            req: &order.CancelOrderRequest{
                OrderNo: "NOT_EXIST_ORDER",
                Reason:  "test cancel",
            },
            wantErr: model.ErrNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewCancelOrderLogic(context.Background(), ctx)
            resp, err := l.CancelOrder(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.True(t, resp.Success)
            }
        })
    }

    // Cleanup
    _ = ctx.OrdersModel.Delete(context.Background(), uint64(orderId))
}