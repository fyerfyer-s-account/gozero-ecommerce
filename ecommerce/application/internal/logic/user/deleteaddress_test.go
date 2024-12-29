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

func TestDeleteAddressLogic_DeleteAddress(t *testing.T) {
	tests := []struct {
		name    string
		userId  int64
		req     *types.DeleteAddressReq
		mock    func(mockUser *User)
		wantErr error
	}{
		{
			name:   "successful deletion",
			userId: 12345,
			req: &types.DeleteAddressReq{
				Id: 67890,
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().DeleteAddress(
					mock.Anything,
					&user.DeleteAddressRequest{
						UserId:    12345,
						AddressId: 67890,
					},
				).Return(&user.DeleteAddressResponse{}, nil)
			},
			wantErr: nil,
		},
		{
			name:   "address not found",
			userId: 12345,
			req: &types.DeleteAddressReq{
				Id: 99999,
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().DeleteAddress(
					mock.Anything,
					&user.DeleteAddressRequest{
						UserId:    12345,
						AddressId: 99999,
					},
				).Return(nil, errors.New("address not found"))
			},
			wantErr: errors.New("address not found"),
		},
		{
			name:   "rpc error",
			userId: 12345,
			req: &types.DeleteAddressReq{
				Id: 67890,
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().DeleteAddress(
					mock.Anything,
					&user.DeleteAddressRequest{
						UserId:    12345,
						AddressId: 67890,
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
			logic := NewDeleteAddressLogic(ctx, svcCtx)

			// Execute DeleteAddress
			err := logic.DeleteAddress(tt.req)

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
