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

func TestConfirmOrderLogic_ConfirmOrder(t *testing.T) {
    configFile := flag.String("f", "../../etc/order.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test orders
    paidOrder := &model.Orders{
        OrderNo:        "TEST_ORDER_001",
        UserId:        1,
        TotalAmount:   100,
        PayAmount:     100,
        FreightAmount: 0,
        Status:        2, // Paid status
        Address:       "Test Address",
        Receiver:      "Test User",
        Phone:         "1234567890",
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    unpaidOrder := &model.Orders{
        OrderNo:        "TEST_ORDER_002",
        UserId:        1,
        TotalAmount:   100,
        PayAmount:     100,
        FreightAmount: 0,
        Status:        1, // Unpaid status
        Address:       "Test Address",
        Receiver:      "Test User",
        Phone:         "1234567890",
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }
    
    result1, _ := ctx.OrdersModel.Insert(context.Background(), paidOrder)
    result2, _ := ctx.OrdersModel.Insert(context.Background(), unpaidOrder)
    paidOrderId, _ := result1.LastInsertId()
    unpaidOrderId, _ := result2.LastInsertId()

    tests := []struct {
        name    string
        req     *order.ConfirmOrderRequest
        wantErr error
    }{
        {
            name: "confirm paid order successfully",
            req: &order.ConfirmOrderRequest{
                OrderNo: "TEST_ORDER_001",
            },
            wantErr: nil,
        },
        {
            name: "invalid order number",
            req: &order.ConfirmOrderRequest{
                OrderNo: "",
            },
            wantErr: zeroerr.ErrOrderNoEmpty,
        },
        {
            name: "order not found",
            req: &order.ConfirmOrderRequest{
                OrderNo: "NOT_EXIST_ORDER",
            },
            wantErr: model.ErrNotFound,
        },
        {
            name: "invalid order status",
            req: &order.ConfirmOrderRequest{
                OrderNo: "TEST_ORDER_002",
            },
            wantErr: zeroerr.ErrOrderStatusNotAllowed,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewConfirmOrderLogic(context.Background(), ctx)
            resp, err := l.ConfirmOrder(tt.req)

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

    _ = ctx.OrdersModel.Delete(context.Background(), uint64(paidOrderId))
    _ = ctx.OrdersModel.Delete(context.Background(), uint64(unpaidOrderId))
}