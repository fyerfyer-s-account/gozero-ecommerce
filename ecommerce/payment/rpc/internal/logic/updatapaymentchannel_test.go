package logic

import (
	"context"
	"encoding/json"
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

func TestUpdatePaymentChannelLogic_UpdatePaymentChannel(t *testing.T) {
    configFile := flag.String("f", "../../etc/payment.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test channel
    testChannel := &model.PaymentChannels{
        Name:      "Test Channel",
        Channel:   1,
        Config:    `{"appId":"test"}`,
        Status:    1,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    result, err := ctx.PaymentChannelsModel.Insert(context.Background(), testChannel)
    assert.NoError(t, err)
    channelId, err := result.LastInsertId()
    assert.NoError(t, err)

    defer func() {
        err := ctx.PaymentChannelsModel.Delete(context.Background(), uint64(channelId))
        assert.NoError(t, err)
    }()

    tests := []struct {
        name    string
        req     *payment.UpdatePaymentChannelRequest
        wantErr error
        check   func(t *testing.T, channel *model.PaymentChannels)
    }{
        {
            name: "Valid Update All Fields",
            req: &payment.UpdatePaymentChannelRequest{
                Id:     channelId,
                Name:   "Updated Channel",
                Config: `{"appId":"new","key":"value"}`,
                Status: 2,
            },
            wantErr: nil,
            check: func(t *testing.T, channel *model.PaymentChannels) {
                assert.Equal(t, "Updated Channel", channel.Name)
                assert.Equal(t, 
					normalizeJSON(`{"appId":"new","key":"value"}`), 
					normalizeJSON(channel.Config),
				)
                assert.Equal(t, int64(2), channel.Status)
            },
        },
        {
            name: "Update Name Only",
            req: &payment.UpdatePaymentChannelRequest{
                Id:   channelId,
                Name: "New Name",
            },
            wantErr: nil,
            check: func(t *testing.T, channel *model.PaymentChannels) {
                assert.Equal(t, "New Name", channel.Name)
            },
        },
        {
            name: "Invalid Channel ID",
            req: &payment.UpdatePaymentChannelRequest{
                Id:   0,
                Name: "Test",
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "Invalid Config JSON",
            req: &payment.UpdatePaymentChannelRequest{
                Id:     channelId,
                Config: `invalid json`,
            },
            wantErr: zeroerr.ErrInvalidChannelConfig,
        },
        {
            name: "Non-existent Channel",
            req: &payment.UpdatePaymentChannelRequest{
                Id:   99999,
                Name: "Test",
            },
            wantErr: zeroerr.ErrChannelNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewUpdatePaymentChannelLogic(context.Background(), ctx)
            resp, err := l.UpdatePaymentChannel(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.True(t, resp.Success)

                if tt.check != nil {
                    channel, err := ctx.PaymentChannelsModel.FindOne(context.Background(), uint64(tt.req.Id))
                    assert.NoError(t, err)
                    tt.check(t, channel)
                }
            }
        })
    }
}

func normalizeJSON(s string) string {
    var obj interface{}
    if err := json.Unmarshal([]byte(s), &obj); err != nil {
        return s
    }
    normalized, _ := json.Marshal(obj)
    return string(normalized)
}
