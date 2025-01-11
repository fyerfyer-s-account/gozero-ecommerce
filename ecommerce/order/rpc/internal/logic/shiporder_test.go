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

func TestShipOrderLogic_ShipOrder(t *testing.T) {
    configFile := flag.String("f", "../../etc/order.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test orders
    confirmedOrder := &model.Orders{
        OrderNo:      "TEST_ORDER_001",
        UserId:      1,
        TotalAmount: 100,
        PayAmount:   100,
        Status:      2, // Confirmed
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    unpaidOrder := &model.Orders{
        OrderNo:      "TEST_ORDER_002",
        UserId:      1,
        TotalAmount: 100,
        PayAmount:   100,
        Status:      1, // Unpaid
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    result1, _ := ctx.OrdersModel.Insert(context.Background(), confirmedOrder)
    result2, _ := ctx.OrdersModel.Insert(context.Background(), unpaidOrder)
    confirmedOrderId, _ := result1.LastInsertId()
    unpaidOrderId, _ := result2.LastInsertId()

    tests := []struct {
        name    string
        req     *order.ShipOrderRequest
        wantErr error
    }{
        {
            name: "ship order successfully",
            req: &order.ShipOrderRequest{
                OrderNo:    "TEST_ORDER_001",
                ShippingNo: "SF123456",
                Company:    "SF Express",
            },
            wantErr: nil,
        },
        {
            name: "empty order number",
            req: &order.ShipOrderRequest{
                OrderNo:    "",
                ShippingNo: "SF123456",
                Company:    "SF Express",
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "empty shipping number",
            req: &order.ShipOrderRequest{
                OrderNo:    "TEST_ORDER_001",
                ShippingNo: "",
                Company:    "SF Express",
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "order not found",
            req: &order.ShipOrderRequest{
                OrderNo:    "NOT_EXIST_ORDER",
                ShippingNo: "SF123456",
                Company:    "SF Express",
            },
            wantErr: model.ErrNotFound,
        },
        {
            name: "invalid order status",
            req: &order.ShipOrderRequest{
                OrderNo:    "TEST_ORDER_002",
                ShippingNo: "SF123456",
                Company:    "SF Express",
            },
            wantErr: zeroerr.ErrOrderStatusNotAllowed,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewShipOrderLogic(context.Background(), ctx)
            resp, err := l.ShipOrder(tt.req)

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
    _ = ctx.OrdersModel.Delete(context.Background(), uint64(confirmedOrderId))
    _ = ctx.OrdersModel.Delete(context.Background(), uint64(unpaidOrderId))
}