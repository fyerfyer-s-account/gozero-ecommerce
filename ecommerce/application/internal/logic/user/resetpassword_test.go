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

func TestResetPasswordLogic_ResetPassword(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.ResetPasswordReq
		mock    func(mockUser *User)
		wantErr error
	}{
		{
			name: "successful password reset",
			req: &types.ResetPasswordReq{
				Phone:    "1234567890",
				Code:     "123456",
				Password: "newPassword123",
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().ResetPassword(
					mock.Anything,
					&user.ResetPasswordRequest{
						Phone:       "1234567890",
						VerifyCode:  "123456",
						NewPassword: "newPassword123",
					},
				).Return(&user.ResetPasswordResponse{}, nil)
			},
			wantErr: nil,
		},
		{
			name: "invalid verification code",
			req: &types.ResetPasswordReq{
				Phone:    "1234567890",
				Code:     "000000",
				Password: "newPassword123",
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().ResetPassword(
					mock.Anything,
					&user.ResetPasswordRequest{
						Phone:       "1234567890",
						VerifyCode:  "000000",
						NewPassword: "newPassword123",
					},
				).Return(nil, errors.New("invalid verification code"))
			},
			wantErr: errors.New("invalid verification code"),
		},
		{
			name: "phone not found",
			req: &types.ResetPasswordReq{
				Phone:    "9999999999",
				Code:     "123456",
				Password: "newPassword123",
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().ResetPassword(
					mock.Anything,
					&user.ResetPasswordRequest{
						Phone:       "9999999999",
						VerifyCode:  "123456",
						NewPassword: "newPassword123",
					},
				).Return(nil, errors.New("phone number not found"))
			},
			wantErr: errors.New("phone number not found"),
		},
		{
			name: "rpc error",
			req: &types.ResetPasswordReq{
				Phone:    "1234567890",
				Code:     "123456",
				Password: "newPassword123",
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().ResetPassword(
					mock.Anything,
					&user.ResetPasswordRequest{
						Phone:       "1234567890",
						VerifyCode:  "123456",
						NewPassword: "newPassword123",
					},
				).Return(nil, errors.New("rpc error"))
			},
			wantErr: errors.New("rpc error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := NewUser(t)
			tt.mock(mockUser)

			svcCtx := &svc.ServiceContext{
				UserRpc: mockUser,
			}

			logic := NewResetPasswordLogic(context.Background(), svcCtx)
			err := logic.ResetPassword(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
