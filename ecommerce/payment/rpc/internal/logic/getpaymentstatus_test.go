package logic

import (
    "context"
    "database/sql"
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

func TestGetPaymentStatusLogic_GetPaymentStatus(t *testing.T) {
    configFile := flag.String("f", "../../etc/payment.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test payments with different statuses
    testCases := []struct {
        paymentNo   string
        status      int64
        channelData string
    }{
        {"PAY_PENDING_001", 1, ""},
        {"PAY_PROCESSING_001", 2, `{"transactionId":"123"}`},
        {"PAY_SUCCESS_001", 3, `{"transactionId":"456","success":true}`},
    }

    var cleanupFuncs []func()
    for _, tc := range testCases {
        payment := &model.PaymentOrders{
            PaymentNo: tc.paymentNo,
            OrderNo:   "TEST_ORDER_001",
            UserId:    1,
            Amount:    100.00,
            Channel:   2,
            Status:    tc.status,
            ChannelData: sql.NullString{
                String: tc.channelData,
                Valid:  tc.channelData != "",
            },
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        }

        result, err := ctx.PaymentOrdersModel.Insert(context.Background(), payment)
        assert.NoError(t, err)
        paymentId, err := result.LastInsertId()
        assert.NoError(t, err)

        cleanupFunc := func(id uint64) func() {
            return func() {
                err := ctx.PaymentOrdersModel.Delete(context.Background(), id)
                assert.NoError(t, err)
            }
        }(uint64(paymentId))
        cleanupFuncs = append(cleanupFuncs, cleanupFunc)
    }

    // Cleanup after all tests
    defer func() {
        for _, cleanup := range cleanupFuncs {
            cleanup()
        }
    }()

    compareJSON := func(t *testing.T, expected, actual string) bool {
        if expected == "" && actual == "" {
            return true
        }
        var expectedMap, actualMap map[string]interface{}
        err := json.Unmarshal([]byte(expected), &expectedMap)
        assert.NoError(t, err)
        err = json.Unmarshal([]byte(actual), &actualMap)
        assert.NoError(t, err)
        return assert.Equal(t, expectedMap, actualMap)
    }

    tests := []struct {
        name          string
        req           *payment.GetPaymentStatusRequest
        wantStatus    int32
        wantData      string
        wantErr       error
    }{
        {
            name:          "Pending Payment",
            req:           &payment.GetPaymentStatusRequest{PaymentNo: "PAY_PENDING_001"},
            wantStatus:    1,
            wantData:      "",
            wantErr:       nil,
        },
        {
            name:          "Processing Payment",
            req:           &payment.GetPaymentStatusRequest{PaymentNo: "PAY_PROCESSING_001"},
            wantStatus:    2,
            wantData:      `{"transactionId":"123"}`,
            wantErr:       nil,
        },
        {
            name:          "Successful Payment",
            req:           &payment.GetPaymentStatusRequest{PaymentNo: "PAY_SUCCESS_001"},
            wantStatus:    3,
            wantData:      `{"transactionId":"456","success":true}`,
            wantErr:       nil,
        },
        {
            name:          "Empty Payment Number",
            req:           &payment.GetPaymentStatusRequest{PaymentNo: ""},
            wantStatus:    0,
            wantData:      "",
            wantErr:       zeroerr.ErrInvalidParam,
        },
        {
            name:          "Non-existent Payment",
            req:           &payment.GetPaymentStatusRequest{PaymentNo: "NON_EXISTENT"},
            wantStatus:    0,
            wantData:      "",
            wantErr:       zeroerr.ErrPaymentNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewGetPaymentStatusLogic(context.Background(), ctx)
            resp, err := l.GetPaymentStatus(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.Equal(t, tt.wantStatus, resp.Status)
                compareJSON(t, tt.wantData, resp.ChannelData)
            }
        })
    }
}
