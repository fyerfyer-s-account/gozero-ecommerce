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

func TestAddAddressLogic_AddAddress(t *testing.T) {
	tests := []struct {
		name    string
		userId  int64
		req     *types.AddressReq
		mock    func(mockUser *User)
		want    *types.Address
		wantErr error
	}{
		{
			name:   "successful address addition",
			userId: 12345,
			req: &types.AddressReq{
				ReceiverName:  "John Doe",
				ReceiverPhone: "1234567890",
				Province:      "TestProvince",
				City:          "TestCity",
				District:      "TestDistrict",
				DetailAddress: "123 Test Street",
				IsDefault:     true,
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().AddAddress(
					mock.Anything,
					&user.AddAddressRequest{
						UserId:        12345,
						ReceiverName:  "John Doe",
						ReceiverPhone: "1234567890",
						Province:      "TestProvince",
						City:          "TestCity",
						District:      "TestDistrict",
						DetailAddress: "123 Test Street",
						IsDefault:     true,
					},
				).Return(&user.AddAddressResponse{
					AddressId: 67890,
				}, nil)
			},
			want: &types.Address{
				Id:            67890,
				ReceiverName:  "John Doe",
				ReceiverPhone: "1234567890",
				Province:      "TestProvince",
				City:          "TestCity",
				District:      "TestDistrict",
				DetailAddress: "123 Test Street",
				IsDefault:     true,
			},
			wantErr: nil,
		},
		{
			name:   "rpc error",
			userId: 12345,
			req: &types.AddressReq{
				ReceiverName:  "John Doe",
				ReceiverPhone: "1234567890",
				Province:      "TestProvince",
				City:          "TestCity",
				District:      "TestDistrict",
				DetailAddress: "123 Test Street",
				IsDefault:     false,
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().AddAddress(
					mock.Anything,
					&user.AddAddressRequest{
						UserId:        12345,
						ReceiverName:  "John Doe",
						ReceiverPhone: "1234567890",
						Province:      "TestProvince",
						City:          "TestCity",
						District:      "TestDistrict",
						DetailAddress: "123 Test Street",
						IsDefault:     false,
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
			logic := NewAddAddressLogic(ctx, svcCtx)

			// Execute AddAddress
			got, err := logic.AddAddress(tt.req)

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
