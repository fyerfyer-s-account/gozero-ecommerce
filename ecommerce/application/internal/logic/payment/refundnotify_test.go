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

func TestRefundNotifyLogic_RefundNotify(t *testing.T) {
    tests := []struct {
        name     string
        req      *types.RefundNotifyReq
        mock     func(mockPayment *Payment)
        wantResp *types.RefundNotifyResp
        wantErr  error
    }{
        {
            name: "successful refund notification",
            req: &types.RefundNotifyReq{
                PaymentType: 1,
                RefundNo:    "REF_001",
                Data:        `{"refund_id":"REF_001","status":"SUCCESS"}`,
            },
            mock: func(mockPayment *Payment) {
                mockPayment.EXPECT().RefundNotify(
                    mock.Anything,
                    &payment.RefundNotifyRequest{
                        Channel:    1,
                        NotifyData: `{"refund_id":"REF_001","status":"SUCCESS"}`,
                    },
                ).Return(&payment.RefundNotifyResponse{
                    ReturnMsg: "success",
                }, nil)
            },
            wantResp: &types.RefundNotifyResp{
                Code:    200,
                Message: "success",
            },
            wantErr: nil,
        },
        {
            name: "invalid parameters",
            req: &types.RefundNotifyReq{
                PaymentType: 0,
                RefundNo:    "",
                Data:        "",
            },
            mock:     func(mockPayment *Payment) {},
            wantResp: nil,
            wantErr:  zeroerr.ErrInvalidParameter,
        },
        {
            name: "rpc error",
            req: &types.RefundNotifyReq{
                PaymentType: 1,
                RefundNo:    "REF_002",
                Data:        `{"refund_id":"REF_002"}`,
            },
            mock: func(mockPayment *Payment) {
                mockPayment.EXPECT().RefundNotify(
                    mock.Anything,
                    mock.Anything,
                ).Return(nil, errors.New("rpc error"))
            },
            wantResp: &types.RefundNotifyResp{
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

            logic := NewRefundNotifyLogic(context.Background(), svcCtx)
            resp, err := logic.RefundNotify(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantResp, resp)
            }
        })
    }
}