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

func TestUpdateProfileLogic_UpdateProfile(t *testing.T) {
	tests := []struct {
		name    string
		userId  int64
		req     *types.UpdateProfileReq
		mock    func(mockUser *User)
		wantErr error
	}{
		{
			name:   "successful profile update",
			userId: 12345,
			req: &types.UpdateProfileReq{
				Nickname: "New Name",
				Avatar:   "new-avatar.jpg",
				Gender:   "male",
				Phone:    "1234567890",
				Email:    "new@example.com",
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().UpdateUserInfo(
					mock.Anything,
					&user.UpdateUserInfoRequest{
						UserId:   12345,
						Nickname: "New Name",
						Avatar:   "new-avatar.jpg",
						Gender:   "male",
						Phone:    "1234567890",
						Email:    "new@example.com",
					},
				).Return(&user.UpdateUserInfoResponse{}, nil)
			},
			wantErr: nil,
		},
		{
			name:   "rpc error",
			userId: 12345,
			req: &types.UpdateProfileReq{
				Nickname: "New Name",
				Avatar:   "new-avatar.jpg",
				Gender:   "male",
				Phone:    "1234567890",
				Email:    "new@example.com",
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().UpdateUserInfo(
					mock.Anything,
					&user.UpdateUserInfoRequest{
						UserId:   12345,
						Nickname: "New Name",
						Avatar:   "new-avatar.jpg",
						Gender:   "male",
						Phone:    "1234567890",
						Email:    "new@example.com",
					},
				).Return(nil, errors.New("rpc error"))
			},
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
			logic := NewUpdateProfileLogic(ctx, svcCtx)

			// Execute UpdateProfile
			err := logic.UpdateProfile(tt.req)

			// Assert results
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
