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

func TestGetRefundLogic_GetRefund(t *testing.T) {
    configFile := flag.String("f", "../../etc/payment.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    testRefund := &model.RefundOrders{
        RefundNo:   "TEST_REF_001",
        PaymentNo:  "TEST_PAY_001",
        OrderNo:    "TEST_ORDER_001",
        UserId:     1,
        Amount:     50.00,
        Reason:     "Test refund",
        Status:     1,
        NotifyUrl:  sql.NullString{String: "http://example.com/notify", Valid: true},
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    result, err := ctx.RefundOrdersModel.Insert(context.Background(), testRefund)
    assert.NoError(t, err)
    refundId, err := result.LastInsertId()
    assert.NoError(t, err)

    defer func() {
        err := ctx.RefundOrdersModel.Delete(context.Background(), uint64(refundId))
        assert.NoError(t, err)
    }()

    tests := []struct {
        name    string
        req     *payment.GetRefundRequest
        wantErr error
    }{
        {
            name: "Valid Refund",
            req: &payment.GetRefundRequest{
                RefundNo: "TEST_REF_001",
            },
            wantErr: nil,
        },
        {
            name: "Empty Refund Number",
            req: &payment.GetRefundRequest{
                RefundNo: "",
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "Not Found Refund",
            req: &payment.GetRefundRequest{
                RefundNo: "NON_EXISTENT",
            },
            wantErr: zeroerr.ErrRefundNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewGetRefundLogic(context.Background(), ctx)
            resp, err := l.GetRefund(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.Equal(t, testRefund.RefundNo, resp.Refund.RefundNo)
                assert.Equal(t, testRefund.OrderNo, resp.Refund.OrderNo)
                assert.Equal(t, testRefund.Amount, resp.Refund.Amount)
                assert.Equal(t, int32(testRefund.Status), resp.Refund.Status)
                assert.Equal(t, testRefund.NotifyUrl.String, resp.Refund.NotifyUrl)
            }
        })
    }
}