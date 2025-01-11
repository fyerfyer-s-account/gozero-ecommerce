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

func TestPayOrderLogic_PayOrder(t *testing.T) {
    configFile := flag.String("f", "../../etc/order.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test orders
    unpaidOrder := &model.Orders{
        OrderNo:      "TEST_ORDER_001",
        UserId:      1,
        TotalAmount: 100,
        PayAmount:   100,
        Status:      1, // Unpaid
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    paidOrder := &model.Orders{
        OrderNo:      "TEST_ORDER_002",
        UserId:      1,
        TotalAmount: 100,
        PayAmount:   100,
        Status:      2, // Paid
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    result1, _ := ctx.OrdersModel.Insert(context.Background(), unpaidOrder)
    result2, _ := ctx.OrdersModel.Insert(context.Background(), paidOrder)
    unpaidOrderId, _ := result1.LastInsertId()
    paidOrderId, _ := result2.LastInsertId()

    tests := []struct {
        name    string
        req     *order.PayOrderRequest
        wantErr error
    }{
        {
            name: "pay order successfully",
            req: &order.PayOrderRequest{
                OrderNo:       "TEST_ORDER_001",
                PaymentMethod: 1,
            },
            wantErr: nil,
        },
        {
            name: "invalid order number",
            req: &order.PayOrderRequest{
                OrderNo:       "",
                PaymentMethod: 1,
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "invalid payment method",
            req: &order.PayOrderRequest{
                OrderNo:       "TEST_ORDER_001",
                PaymentMethod: 0,
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "order not found",
            req: &order.PayOrderRequest{
                OrderNo:       "NOT_EXIST_ORDER",
                PaymentMethod: 1,
            },
            wantErr: model.ErrNotFound,
        },
        {
            name: "order already paid",
            req: &order.PayOrderRequest{
                OrderNo:       "TEST_ORDER_002",
                PaymentMethod: 1,
            },
            wantErr: zeroerr.ErrOrderStatusNotAllowed,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewPayOrderLogic(context.Background(), ctx)
            resp, err := l.PayOrder(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.NotEmpty(t, resp.PaymentNo)
                assert.NotEmpty(t, resp.PayUrl)
            }
        })
    }

    // Cleanup
    _ = ctx.OrdersModel.Delete(context.Background(), uint64(unpaidOrderId))
    _ = ctx.OrdersModel.Delete(context.Background(), uint64(paidOrderId))
}