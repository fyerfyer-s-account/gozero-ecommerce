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

func TestPaymentNotifyLogic_PaymentNotify(t *testing.T) {
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
        Channel:   2, // Alipay
        Status:    2, // Processing
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    result, err := ctx.PaymentOrdersModel.Insert(context.Background(), testPayment)
    assert.NoError(t, err)
    paymentId, err := result.LastInsertId()
    assert.NoError(t, err)

    defer func() {
        err := ctx.PaymentOrdersModel.Delete(context.Background(), uint64(paymentId))
        assert.NoError(t, err)
    }()

    tests := []struct {
        name      string
        req       *payment.PaymentNotifyRequest
        wantMsg   string
        wantErr   error
        checkFunc func(t *testing.T, ctx *svc.ServiceContext)
    }{
        {
            name: "Valid Alipay Success",
            req: &payment.PaymentNotifyRequest{
                Channel: 2,
                NotifyData: `{"out_trade_no":"TEST_PAY_001","trade_status":"TRADE_SUCCESS"}`,
            },
            wantMsg: "success",
            wantErr: nil,
            checkFunc: func(t *testing.T, ctx *svc.ServiceContext) {
                payment, err := ctx.PaymentOrdersModel.FindOneByPaymentNo(context.Background(), "TEST_PAY_001")
                assert.NoError(t, err)
                assert.Equal(t, int64(3), payment.Status) // Paid
                assert.True(t, payment.PayTime.Valid)
            },
        },
        {
            name: "Invalid Data Format",
            req: &payment.PaymentNotifyRequest{
                Channel: 2,
                NotifyData: `invalid json`,
            },
            wantMsg: "",
            wantErr: zeroerr.ErrInvalidNotifyData,
        },
        {
            name: "Unknown Channel",
            req: &payment.PaymentNotifyRequest{
                Channel: 99,
                NotifyData: `{"out_trade_no":"TEST_PAY_001"}`,
            },
            wantMsg: "",
            wantErr: zeroerr.ErrUnsupportedChannel,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewPaymentNotifyLogic(context.Background(), ctx)
            resp, err := l.PaymentNotify(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.Equal(t, tt.wantMsg, resp.ReturnMsg)

                if tt.checkFunc != nil {
                    tt.checkFunc(t, ctx)
                }
            }
        })
    }
}