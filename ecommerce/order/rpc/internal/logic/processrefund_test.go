package logic

import (
    "context"
    "database/sql"
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

func TestProcessRefundLogic_ProcessRefund(t *testing.T) {
    configFile := flag.String("f", "../../etc/order.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test data
    testOrder := &model.Orders{
        OrderNo:      "TEST_ORDER_001",
        UserId:      1,
        TotalAmount: 100,
        PayAmount:   100,
        Status:      6, // Refunding
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    orderResult, _ := ctx.OrdersModel.Insert(context.Background(), testOrder)
    orderId, _ := orderResult.LastInsertId()

    testPayment := &model.OrderPayments{
        OrderId:       uint64(orderId),
        PaymentNo:     "PAY_TEST_001",
        PaymentMethod: 1,
        Amount:        100,
        Status:        1, // Paid
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }
    
    ctx.OrderPaymentsModel.Insert(context.Background(), testPayment)

    testRefund := &model.OrderRefunds{
        OrderId:     uint64(orderId),
        RefundNo:    "RF_TEST_001",
        Amount:      50,
        Reason:      "test refund",
        Status:      0, // Pending
        Description: sql.NullString{String: "test", Valid: true},
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    processedRefund := &model.OrderRefunds{
        OrderId:     uint64(orderId),
        RefundNo:    "RF_TEST_002",
        Amount:      50,
        Status:      1, // Already processed
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    refundResult1, _ := ctx.OrderRefundsModel.Insert(context.Background(), testRefund)
    refundResult2, _ := ctx.OrderRefundsModel.Insert(context.Background(), processedRefund)
    refundId1, _ := refundResult1.LastInsertId()
    refundId2, _ := refundResult2.LastInsertId()

    tests := []struct {
        name    string
        req     *order.ProcessRefundRequest
        wantErr error
    }{
        {
            name: "approve refund successfully",
            req: &order.ProcessRefundRequest{
                RefundNo: "RF_TEST_001",
                Agree:    true,
                Reply:    "approved",
            },
            wantErr: nil,
        },
        {
            name: "empty refund number",
            req: &order.ProcessRefundRequest{
                RefundNo: "",
                Agree:    true,
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "refund not found",
            req: &order.ProcessRefundRequest{
                RefundNo: "NOT_EXIST_REFUND",
                Agree:    true,
            },
            wantErr: model.ErrNotFound,
        },
        {
            name: "already processed refund",
            req: &order.ProcessRefundRequest{
                RefundNo: "RF_TEST_002",
                Agree:    true,
            },
            wantErr: zeroerr.ErrRefundStatusInvalid,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewProcessRefundLogic(context.Background(), ctx)
            resp, err := l.ProcessRefund(tt.req)

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
    _ = ctx.OrderRefundsModel.Delete(context.Background(), uint64(refundId1))
    _ = ctx.OrderRefundsModel.Delete(context.Background(), uint64(refundId2))
}