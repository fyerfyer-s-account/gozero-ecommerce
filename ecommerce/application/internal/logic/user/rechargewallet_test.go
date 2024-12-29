package user

import (
	"context"
	"errors"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRechargeWalletLogic_RechargeWallet(t *testing.T) {
	tests := []struct {
		name    string
		userId  int64
		req     *types.RechargeReq
		mock    func(mockUser *User)
		wantErr error
	}{
		{
			name:   "successful recharge with alipay",
			userId: 12345,
			req: &types.RechargeReq{
				Amount:      100.50,
				PaymentType: 1, // alipay
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().RechargeWallet(
					mock.Anything,
					&user.RechargeWalletRequest{
						UserId:  12345,
						Amount:  100.50,
						Channel: "alipay",
					},
				).Return(&user.RechargeWalletResponse{}, nil)
			},
			wantErr: nil,
		},
		{
			name:   "successful recharge with wechat",
			userId: 12345,
			req: &types.RechargeReq{
				Amount:      200.75,
				PaymentType: 2, // wechat
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().RechargeWallet(
					mock.Anything,
					&user.RechargeWalletRequest{
						UserId:  12345,
						Amount:  200.75,
						Channel: "wechat",
					},
				).Return(&user.RechargeWalletResponse{}, nil)
			},
			wantErr: nil,
		},
		{
			name:   "recharge failed - rpc error",
			userId: 12345,
			req: &types.RechargeReq{
				Amount:      100.00,
				PaymentType: 1,
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().RechargeWallet(
					mock.Anything,
					&user.RechargeWalletRequest{
						UserId:  12345,
						Amount:  100.00,
						Channel: "alipay",
					},
				).Return(nil, errors.New("rpc error"))
			},
			wantErr: zeroerr.ErrRechargeWalletFailed,
		},
		{
			name:   "recharge with unknown payment type",
			userId: 12345,
			req: &types.RechargeReq{
				Amount:      100.00,
				PaymentType: 3, // unknown
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().RechargeWallet(
					mock.Anything,
					&user.RechargeWalletRequest{
						UserId:  12345,
						Amount:  100.00,
						Channel: "unknown",
					},
				).Return(nil, errors.New("invalid payment type"))
			},
			wantErr: zeroerr.ErrRechargeWalletFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock user client
			mockUser := NewUser(t)

			// Setup mock expectations
			tt.mock(mockUser)

			// Create service context with mock
			svcCtx := &svc.ServiceContext{
				UserRpc: mockUser,
			}

			// Create context with userId
			ctx := context.WithValue(context.Background(), "userId", tt.userId)

			// Create logic instance
			logic := NewRechargeWalletLogic(ctx, svcCtx)

			// Execute RechargeWallet
			err := logic.RechargeWallet(tt.req)

			// Assert results
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
