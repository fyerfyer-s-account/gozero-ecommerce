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
    "github.com/stretchr/testify/assert"
    "github.com/zeromicro/go-zero/core/conf"
)

func TestListPaymentChannelsLogic_ListPaymentChannels(t *testing.T) {
    configFile := flag.String("f", "../../etc/payment.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test channels
    testChannels := []*model.PaymentChannels{
        {
            Name:      "Test Wechat",
            Channel:   1,
            Config:    `{"appId":"wx1"}`,
            Status:    1,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
        {
            Name:      "Test Alipay",
            Channel:   2,
            Config:    `{"appId":"ali1"}`,
            Status:    1,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
        {
            Name:      "Test Disabled",
            Channel:   3,
            Config:    `{"appId":"test"}`,
            Status:    2,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
    }

    var cleanupFuncs []func()
    for _, ch := range testChannels {
        result, err := ctx.PaymentChannelsModel.Insert(context.Background(), ch)
        assert.NoError(t, err)
        id, err := result.LastInsertId()
        assert.NoError(t, err)

        cleanupFunc := func(id uint64) func() {
            return func() {
                err := ctx.PaymentChannelsModel.Delete(context.Background(), id)
                assert.NoError(t, err)
            }
        }(uint64(id))
        cleanupFuncs = append(cleanupFuncs, cleanupFunc)
    }

    defer func() {
        for _, cleanup := range cleanupFuncs {
            cleanup()
        }
    }()

    tests := []struct {
        name          string
        req           *payment.ListPaymentChannelsRequest
        wantCount     int
        wantStatus    int32
    }{
        {
            name:          "List All Channels",
            req:           &payment.ListPaymentChannelsRequest{},
            wantCount:     3,
            wantStatus:    0,
        },
        {
            name:          "List Enabled Channels",
            req:           &payment.ListPaymentChannelsRequest{Status: 1},
            wantCount:     2,
            wantStatus:    1,
        },
        {
            name:          "List Disabled Channels",
            req:           &payment.ListPaymentChannelsRequest{Status: 2},
            wantCount:     1,
            wantStatus:    2,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewListPaymentChannelsLogic(context.Background(), ctx)
            resp, err := l.ListPaymentChannels(tt.req)

            assert.NoError(t, err)
            assert.NotNil(t, resp)
            assert.Equal(t, tt.wantCount, len(resp.Channels))

            if tt.wantStatus > 0 {
                for _, ch := range resp.Channels {
                    assert.Equal(t, tt.wantStatus, ch.Status)
                }
            }
        })
    }
}