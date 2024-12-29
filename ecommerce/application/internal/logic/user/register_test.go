package user

import (
	"context"
	"errors"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
	"github.com/stretchr/testify/assert"
)

func TestRegisterLogic_Register(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.RegisterReq
		mock    func(mockUser *User)
		want    *types.RegisterResp
		wantErr error
	}{
		{
			name: "successful registration",
			req: &types.RegisterReq{
				Username: "testuser",
				Password: "password123",
				Phone:    "1234567890",
				Email:    "test@example.com",
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().Register(
					context.Background(),
					&user.RegisterRequest{
						Username: "testuser",
						Password: "password123",
						Phone:    "1234567890",
						Email:    "test@example.com",
					},
				).Return(&user.RegisterResponse{
					UserId: 12345,
				}, nil)
			},
			want: &types.RegisterResp{
				UserId: 12345,
			},
			wantErr: nil,
		},
		{
			name: "failed registration - user exists",
			req: &types.RegisterReq{
				Username: "existinguser",
				Password: "password123",
				Phone:    "1234567890",
				Email:    "existing@example.com",
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().Register(
					context.Background(),
					&user.RegisterRequest{
						Username: "existinguser",
						Password: "password123",
						Phone:    "1234567890",
						Email:    "existing@example.com",
					},
				).Return(nil, errors.New("user already exists"))
			},
			want:    nil,
			wantErr: errors.New("user already exists"),
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

			// Create register logic instance
			logic := NewRegisterLogic(context.Background(), svcCtx)

			// Execute register
			got, err := logic.Register(tt.req)

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
