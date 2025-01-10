package logic

import (
	"context"
	"flag"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestCreatePaymentChannelLogic_CreatePaymentChannel(t *testing.T) {
    configFile := flag.String("f", "../../etc/payment.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    
    ctx := svc.NewServiceContext(c)

    tests := []struct {
        name    string
        req     *payment.CreatePaymentChannelRequest
        wantErr error
    }{
        {
            name: "Valid Channel",
            req: &payment.CreatePaymentChannelRequest{
                Name:    "Test Alipay",
                Channel: 2, // Alipay
                Config:  `{"appId":"test","privateKey":"test","publicKey":"test"}`,
            },
            wantErr: nil,
        },
        {
            name: "Invalid Params - Empty Name",
            req: &payment.CreatePaymentChannelRequest{
                Name:    "",
                Channel: 2,
                Config:  `{"appId":"test"}`,
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "Invalid Params - Invalid Channel",
            req: &payment.CreatePaymentChannelRequest{
                Name:    "Test Channel",
                Channel: 0,
                Config:  `{"appId":"test"}`,
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "Invalid Config JSON",
            req: &payment.CreatePaymentChannelRequest{
                Name:    "Test Channel",
                Channel: 2,
                Config:  `{invalid json}`,
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewCreatePaymentChannelLogic(context.Background(), ctx)
            resp, err := l.CreatePaymentChannel(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.NotZero(t, resp.Id)

                err = ctx.PaymentChannelsModel.Delete(context.Background(), uint64(resp.Id))
                assert.NoError(t, err)
            }
        })
    }
}