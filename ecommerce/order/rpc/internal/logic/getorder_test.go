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

func TestGetOrderLogic_GetOrder(t *testing.T) {
    configFile := flag.String("f", "../../etc/order.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test order with items
    testOrder := &model.Orders{
        OrderNo:      "TEST_ORDER_001",
        UserId:      1,
        TotalAmount: 100,
        PayAmount:   100,
        Status:      1,
        Address:     "Test Address",
        Receiver:    "Test User",
        Phone:       "1234567890",
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    result, _ := ctx.OrdersModel.Insert(context.Background(), testOrder)
    orderId, _ := result.LastInsertId()

    orderItem := &model.OrderItems{
        OrderId:     uint64(orderId),
        ProductId:   1,
        SkuId:      1,
        ProductName: "Test Product",
        SkuName:    "Test SKU",
        Price:      50,
        Quantity:   2,
        TotalAmount: 100,
        CreatedAt:  time.Now(),
    }
    
    ctx.OrderItemsModel.Insert(context.Background(), orderItem)

    tests := []struct {
        name    string
        req     *order.GetOrderRequest
        wantErr error
    }{
        {
            name: "get order successfully",
            req: &order.GetOrderRequest{
                OrderNo: "TEST_ORDER_001",
            },
            wantErr: nil,
        },
        {
            name: "empty order number",
            req: &order.GetOrderRequest{
                OrderNo: "",
            },
            wantErr: zeroerr.ErrOrderNoEmpty,
        },
        {
            name: "order not found",
            req: &order.GetOrderRequest{
                OrderNo: "NOT_EXIST_ORDER",
            },
            wantErr: model.ErrNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewGetOrderLogic(context.Background(), ctx)
            resp, err := l.GetOrder(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.Equal(t, "TEST_ORDER_001", resp.Order.OrderNo)
                assert.Equal(t, float64(100), resp.Order.TotalAmount)
                assert.Equal(t, int32(1), resp.Order.Status)
                assert.Len(t, resp.Order.Items, 1)
            }
        })
    }

    // Cleanup
    _ = ctx.OrdersModel.Delete(context.Background(), uint64(orderId))
}