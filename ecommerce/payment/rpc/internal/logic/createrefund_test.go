package logic

import (
    "context"
    "flag"
    "testing"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
    "github.com/stretchr/testify/assert"
    "github.com/zeromicro/go-zero/core/conf"
)

func TestCreateRefundLogic_CreateRefund(t *testing.T) {
    configFile := flag.String("f", "../../etc/payment.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test payment order
    testPayment := &model.PaymentOrders{
        PaymentNo: "TEST_PAY_001",
        OrderNo:   "TEST_ORDER_001",
        UserId:    1,
        Amount:    100.00,
        Channel:   2,
        Status:    3, // Paid
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    result, err := ctx.PaymentOrdersModel.Insert(context.Background(), testPayment)
    assert.NoError(t, err)
    paymentId, err := result.LastInsertId()
    assert.NoError(t, err)

    // Cleanup
    defer func() {
        err := ctx.PaymentOrdersModel.Delete(context.Background(), uint64(paymentId))
        assert.NoError(t, err)
    }()

    tests := []struct {
        name    string
        req     *payment.CreateRefundRequest
        wantErr error
    }{
        {
            name: "Valid Refund",
            req: &payment.CreateRefundRequest{
                PaymentNo:  "TEST_PAY_001",
                OrderNo:    "TEST_ORDER_001",
                Amount:    50.00,
                Reason:    "Test refund",
                NotifyUrl: "http://example.com/notify",
            },
            wantErr: nil,
        },
        {
            name: "Invalid Payment Number",
            req: &payment.CreateRefundRequest{
                PaymentNo:  "",
                OrderNo:    "TEST_ORDER_001",
                Amount:    50.00,
                Reason:    "Test refund",
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "Invalid Amount",
            req: &payment.CreateRefundRequest{
                PaymentNo:  "TEST_PAY_001",
                OrderNo:    "TEST_ORDER_001",
                Amount:    150.00,
                Reason:    "Test refund",
            },
            wantErr: zeroerr.ErrRefundAmountInvalid,
        },
        {
            name: "Payment Not Found",
            req: &payment.CreateRefundRequest{
                PaymentNo:  "NON_EXISTENT",
                OrderNo:    "TEST_ORDER_001",
                Amount:    50.00,
                Reason:    "Test refund",
            },
            wantErr: zeroerr.ErrPaymentNotFound,
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

                // Cleanup created refund
                refunds, err := ctx.RefundOrdersModel.FindByRefundNo(context.Background(), resp.RefundNo)
                assert.NoError(t, err)
                if len(refunds) > 0 {
                    err = ctx.RefundOrdersModel.Delete(context.Background(), refunds[0].Id)
                    assert.NoError(t, err)
                }
            }
        })
    }
}