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

// Helper function to create context with userId
func createContextWithUserId(userId int64) context.Context {
	return context.WithValue(context.Background(), "userId", userId)
}

func TestGetProfileLogic_GetProfile(t *testing.T) {
	tests := []struct {
		name    string
		userId  int64
		mock    func(mockUser *User)
		want    *types.UserInfo
		wantErr error
	}{
		{
			name:   "successful profile retrieval",
			userId: 12345,
			mock: func(mockUser *User) {
				mockUser.EXPECT().GetUserInfo(
					mock.Anything, // Use mock.Anything() instead of specific context
					&user.GetUserInfoRequest{
						UserId: 12345,
					},
				).Return(&user.GetUserInfoResponse{
					UserInfo: &user.UserInfo{
						UserId:        12345,
						Username:      "testuser",
						Nickname:      "Test User",
						Avatar:        "avatar.jpg",
						Phone:         "1234567890",
						Email:         "test@example.com",
						Gender:        "male",
						MemberLevel:   1,
						WalletBalance: 100.50,
						CreatedAt:     1634567890,
					},
				}, nil)
			},
			want: &types.UserInfo{
				Id:          12345,
				Username:    "testuser",
				Nickname:    "Test User",
				Avatar:      "avatar.jpg",
				Phone:       "1234567890",
				Email:       "test@example.com",
				Gender:      "male",
				MemberLevel: 1,
				Balance:     100.50,
				CreatedAt:   1634567890,
			},
			wantErr: nil,
		},
		{
			name:   "rpc error",
			userId: 12345,
			mock: func(mockUser *User) {
				mockUser.EXPECT().GetUserInfo(
					mock.Anything, // Use mock.Anything() instead of specific context
					&user.GetUserInfoRequest{
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

			// Create logic with context containing userId
			ctx := createContextWithUserId(tt.userId)
			logic := NewGetProfileLogic(ctx, svcCtx)

			// Execute GetProfile
			got, err := logic.GetProfile()

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
