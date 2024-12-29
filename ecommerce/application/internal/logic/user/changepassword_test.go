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

func TestChangePasswordLogic_ChangePassword(t *testing.T) {
	tests := []struct {
		name    string
		userId  int64
		req     *types.ChangePasswordReq
		mock    func(mockUser *User)
		wantErr error
	}{
		{
			name:   "successful password change",
			userId: 12345,
			req: &types.ChangePasswordReq{
				OldPassword: "oldPass123",
				NewPassword: "newPass123",
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().ChangePassword(
					mock.Anything,
					&user.ChangePasswordRequest{
						UserId:      12345,
						OldPassword: "oldPass123",
						NewPassword: "newPass123",
					},
				).Return(&user.ChangePasswordResponse{}, nil)
			},
			wantErr: nil,
		},
		{
			name:   "invalid old password",
			userId: 12345,
			req: &types.ChangePasswordReq{
				OldPassword: "wrongPass",
				NewPassword: "newPass123",
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().ChangePassword(
					mock.Anything,
					&user.ChangePasswordRequest{
						UserId:      12345,
						OldPassword: "wrongPass",
						NewPassword: "newPass123",
					},
				).Return(nil, errors.New("invalid old password"))
			},
			wantErr: errors.New("invalid old password"),
		},
		{
			name:   "rpc error",
			userId: 12345,
			req: &types.ChangePasswordReq{
				OldPassword: "oldPass123",
				NewPassword: "newPass123",
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().ChangePassword(
					mock.Anything,
					&user.ChangePasswordRequest{
						UserId:      12345,
						OldPassword: "oldPass123",
						NewPassword: "newPass123",
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
			logic := NewChangePasswordLogic(ctx, svcCtx)

			// Execute ChangePassword
			err := logic.ChangePassword(tt.req)

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
