package payment

import (
	"context"
	"errors"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPaymentNotifyLogic_PaymentNotify(t *testing.T) {
    tests := []struct {
        name       string
        req        *types.PaymentNotifyReq
        mock       func(mockPayment *Payment)
        wantResp   *types.PaymentNotifyResp
        wantErr    error
    }{
        {
            name: "successful notification",
            req: &types.PaymentNotifyReq{
                PaymentType: 1,
                PaymentNo:   "TEST_PAY_001",
                Data:       `{"out_trade_no":"TEST_PAY_001","trade_status":"SUCCESS"}`,
            },
            mock: func(mockPayment *Payment) {
                mockPayment.EXPECT().PaymentNotify(
                    mock.Anything,
                    &payment.PaymentNotifyRequest{
                        Channel:    1,
                        NotifyData: `{"out_trade_no":"TEST_PAY_001","trade_status":"SUCCESS"}`,
                    },
                ).Return(&payment.PaymentNotifyResponse{
                    ReturnMsg: "success",
                }, nil)
            },
            wantResp: &types.PaymentNotifyResp{
                Code:    200,
                Message: "success",
            },
            wantErr: nil,
        },
        {
            name: "invalid payment type",
            req: &types.PaymentNotifyReq{
                PaymentType: 0,
                PaymentNo:   "TEST_PAY_001",
                Data:       "test",
            },
            mock:     func(mockPayment *Payment) {},
            wantResp: nil,
            wantErr:  zeroerr.ErrInvalidParameter,
        },
        {
            name: "rpc error",
            req: &types.PaymentNotifyReq{
                PaymentType: 1,
                PaymentNo:   "TEST_PAY_001",
                Data:       "test",
            },
            mock: func(mockPayment *Payment) {
                mockPayment.EXPECT().PaymentNotify(
                    mock.Anything,
                    mock.Anything,
                ).Return(nil, errors.New("rpc error"))
            },
            wantResp: &types.PaymentNotifyResp{
                Code:    500,
                Message: "rpc error",
            },
            wantErr: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockPayment := NewPayment(t)
            tt.mock(mockPayment)

            svcCtx := &svc.ServiceContext{
                PaymentRpc: mockPayment,
            }

            logic := NewPaymentNotifyLogic(context.Background(), svcCtx)
            resp, err := logic.PaymentNotify(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantResp, resp)
            }
        })
    }
}