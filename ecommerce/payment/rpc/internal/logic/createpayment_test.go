package logic

import (
	"context"
	"flag"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestCreatePaymentLogic_CreatePayment(t *testing.T) {
    configFile := flag.String("f", "../../etc/payment.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test payment channel
    channel := &model.PaymentChannels{
        Name:    "Test Alipay",
        Channel: 2,
        Config:  `{"appId":"test","privateKey":"test"}`,
        Status:  1,
    }
    result, err := ctx.PaymentChannelsModel.Insert(context.Background(), channel)
    assert.NoError(t, err)
    channelId, err := result.LastInsertId()
    assert.NoError(t, err)

    // Cleanup
    defer func() {
        err := ctx.PaymentChannelsModel.Delete(context.Background(), uint64(channelId))
        assert.NoError(t, err)
    }()

    tests := []struct {
        name    string
        req     *payment.CreatePaymentRequest
        wantErr error
    }{
        {
            name: "Valid Payment",
            req: &payment.CreatePaymentRequest{
                OrderNo:    "TEST123",
                UserId:    1,
                Amount:    100.00,
                Channel:   2,
                NotifyUrl: "http://example.com/notify",
            },
            wantErr: nil,
        },
        {
            name: "Invalid Params - Empty OrderNo",
            req: &payment.CreatePaymentRequest{
                OrderNo:    "",
                UserId:    1,
                Amount:    100.00,
                Channel:   2,
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "Invalid Amount - Too Low",
            req: &payment.CreatePaymentRequest{
                OrderNo:    "TEST123",
                UserId:    1,
                Amount:    0.001,
                Channel:   2,
            },
            wantErr: zeroerr.ErrInvalidAmount,
        },
        {
            name: "Invalid Channel",
            req: &payment.CreatePaymentRequest{
                OrderNo:    "TEST123",
                UserId:    1,
                Amount:    100.00,
                Channel:   99,
            },
            wantErr: zeroerr.ErrChannelNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewCreatePaymentLogic(context.Background(), ctx)
            resp, err := l.CreatePayment(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.NotEmpty(t, resp.PaymentNo)
                assert.NotEmpty(t, resp.PayUrl)

                // Cleanup created payment
                payments, err := ctx.PaymentOrdersModel.FindByPaymentNo(context.Background(), resp.PaymentNo)
                assert.NoError(t, err)
                if len(payments) > 0 {
                    err = ctx.PaymentOrdersModel.Delete(context.Background(), payments[0].Id)
                    assert.NoError(t, err)
                }
            }
        })
    }
}