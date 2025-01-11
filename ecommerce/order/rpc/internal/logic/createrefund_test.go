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

func TestCreateRefundLogic_CreateRefund(t *testing.T) {
    configFile := flag.String("f", "../../etc/order.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test order
    testOrder := &model.Orders{
        OrderNo:      "TEST_ORDER_001",
        UserId:      1,
        TotalAmount: 100,
        PayAmount:   100,
        Status:      2, // Paid status
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    result, err := ctx.OrdersModel.Insert(context.Background(), testOrder)
    assert.NoError(t, err)
    orderId, err := result.LastInsertId()
    assert.NoError(t, err)

    tests := []struct {
        name    string
        req     *order.CreateRefundRequest
        wantErr error
    }{
        {
            name: "create refund successfully",
            req: &order.CreateRefundRequest{
                OrderNo:     "TEST_ORDER_001",
                Amount:      50,
                Reason:      "test refund",
                Description: "test description",
                Images:      []string{"image1.jpg", "image2.jpg"},
            },
            wantErr: nil,
        },
        {
            name: "invalid order number",
            req: &order.CreateRefundRequest{
                OrderNo: "",
                Amount:  50,
                Reason: "test refund",
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "invalid amount",
            req: &order.CreateRefundRequest{
                OrderNo: "TEST_ORDER_001",
                Amount:  0,
                Reason: "test refund",
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "exceed refund amount",
            req: &order.CreateRefundRequest{
                OrderNo: "TEST_ORDER_001",
                Amount:  150,
                Reason: "test refund",
            },
            wantErr: zeroerr.ErrRefundExceedAmount,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewCreateRefundLogic(context.Background(), ctx)
            resp, err := l.CreateRefund(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.NotEmpty(t, resp.RefundNo)
            }
        })
    }

    // Cleanup
    _ = ctx.OrdersModel.Delete(context.Background(), uint64(orderId))
}