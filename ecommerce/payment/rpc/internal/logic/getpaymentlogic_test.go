package logic

import (
    "context"
    "database/sql"
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

func TestGetPaymentLogic_GetPayment(t *testing.T) {
    configFile := flag.String("f", "../../etc/payment.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test payment
    testPayment := &model.PaymentOrders{
        PaymentNo: "TEST_PAY_001",
        OrderNo:   "TEST_ORDER_001",
        UserId:    1,
        Amount:    100.00,
        Channel:   2,
        Status:    1,
        NotifyUrl: sql.NullString{String: "http://example.com/notify", Valid: true},
        ReturnUrl: sql.NullString{String: "http://example.com/return", Valid: true},
        ExpireTime: sql.NullTime{
            Time:  time.Now().Add(2 * time.Hour),
            Valid: true,
        },
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
        req     *payment.GetPaymentRequest
        wantErr error
    }{
        {
            name: "Valid Payment",
            req: &payment.GetPaymentRequest{
                PaymentNo: "TEST_PAY_001",
            },
            wantErr: nil,
        },
        {
            name: "Empty Payment Number",
            req: &payment.GetPaymentRequest{
                PaymentNo: "",
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "Not Found Payment",
            req: &payment.GetPaymentRequest{
                PaymentNo: "NON_EXISTENT",
            },
            wantErr: zeroerr.ErrPaymentNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewGetPaymentLogic(context.Background(), ctx)
            resp, err := l.GetPayment(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.Equal(t, testPayment.PaymentNo, resp.Payment.PaymentNo)
                assert.Equal(t, testPayment.OrderNo, resp.Payment.OrderNo)
                assert.Equal(t, testPayment.Amount, resp.Payment.Amount)
                assert.Equal(t, testPayment.Channel, resp.Payment.Channel)
                assert.Equal(t, int64(testPayment.Status), resp.Payment.Status)
            }
        })
    }
}