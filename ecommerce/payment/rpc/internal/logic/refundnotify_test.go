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

func TestRefundNotifyLogic_RefundNotify(t *testing.T) {
    configFile := flag.String("f", "../../etc/payment.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test payment and refund orders
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

    paymentResult, err := ctx.PaymentOrdersModel.Insert(context.Background(), testPayment)
    assert.NoError(t, err)
    paymentId, err := paymentResult.LastInsertId()
    assert.NoError(t, err)

    testRefund := &model.RefundOrders{
        RefundNo:   "TEST_REF_001",
        PaymentNo:  "TEST_PAY_001",
        OrderNo:    "TEST_ORDER_001",
        UserId:     1,
        Amount:     50.00,
        Reason:     "Test refund",
        Status:     2, // Processing
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    refundResult, err := ctx.RefundOrdersModel.Insert(context.Background(), testRefund)
    assert.NoError(t, err)
    refundId, err := refundResult.LastInsertId()
    assert.NoError(t, err)

    // Cleanup
    defer func() {
        _ = ctx.PaymentOrdersModel.Delete(context.Background(), uint64(paymentId))
        _ = ctx.RefundOrdersModel.Delete(context.Background(), uint64(refundId))
    }()

    tests := []struct {
        name    string
        req     *payment.RefundNotifyRequest
        wantErr error
        check   func(t *testing.T, ctx *svc.ServiceContext)
    }{
        {
            name: "Valid Alipay Success",
            req: &payment.RefundNotifyRequest{
                Channel:    2,
                NotifyData: `{"out_refund_no":"TEST_REF_001","refund_status":"REFUND_SUCCESS"}`,
            },
            wantErr: nil,
            check: func(t *testing.T, ctx *svc.ServiceContext) {
                refund, err := ctx.RefundOrdersModel.FindOneByRefundNo(context.Background(), "TEST_REF_001")
                assert.NoError(t, err)
                assert.Equal(t, int64(3), refund.Status) // Refunded

                payment, err := ctx.PaymentOrdersModel.FindOneByPaymentNo(context.Background(), "TEST_PAY_001")
                assert.NoError(t, err)
                assert.Equal(t, int64(4), payment.Status) // Refunded
            },
        },
        {
            name: "Invalid Channel",
            req: &payment.RefundNotifyRequest{
                Channel:    99,
                NotifyData: `{"refund_no":"TEST_REF_001"}`,
            },
            wantErr: zeroerr.ErrUnsupportedChannel,
        },
        {
            name: "Invalid Data Format",
            req: &payment.RefundNotifyRequest{
                Channel:    2,
                NotifyData: `invalid json`,
            },
            wantErr: zeroerr.ErrInvalidNotifyData,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewRefundNotifyLogic(context.Background(), ctx)
            resp, err := l.RefundNotify(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.Equal(t, "success", resp.ReturnMsg)

                if tt.check != nil {
                    tt.check(t, ctx)
                }
            }
        })
    }
}