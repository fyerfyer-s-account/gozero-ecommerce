package user

import (
	"context"
	"errors"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetWalletLogic_GetWallet(t *testing.T) {
	tests := []struct {
		name    string
		userId  int64
		mock    func(mockUser *User)
		want    *types.WalletDetail
		wantErr error
	}{
		{
			name:   "successful wallet retrieval",
			userId: 12345,
			mock: func(mockUser *User) {
				mockUser.EXPECT().GetWallet(
					mock.Anything,
					&user.GetWalletRequest{
						UserId: 12345,
					},
				).Return(&user.GetWalletResponse{
					Balance:      1000.50,
					Status:       1,
					FreezeAmount: 100.00,
				}, nil)
			},
			want: &types.WalletDetail{
				Balance:      1000.50,
				Status:       1,
				FrozenAmount: 100.00,
			},
			wantErr: nil,
		},
		{
			name:   "rpc error",
			userId: 12345,
			mock: func(mockUser *User) {
				mockUser.EXPECT().GetWallet(
					mock.Anything,
					&user.GetWalletRequest{
						UserId: 12345,
					},
				).Return(nil, errors.New("rpc error"))
			},
			want:    nil,
			wantErr: errors.New("rpc error"),
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
			logic := NewGetWalletLogic(ctx, svcCtx)

			// Execute GetWallet
			got, err := logic.GetWallet()

			// Assert results
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
