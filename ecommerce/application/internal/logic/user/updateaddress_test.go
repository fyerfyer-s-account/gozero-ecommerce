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

func TestUpdateAddressLogic_UpdateAddress(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.AddressReq
		mock    func(mockUser *User)
		wantErr error
	}{
		{
			name: "successful address update",
			req: &types.AddressReq{
				AddressId:     1,
				ReceiverName:  "John Doe",
				ReceiverPhone: "1234567890",
				Province:      "TestProvince",
				City:          "TestCity",
				District:      "TestDistrict",
				DetailAddress: "123 Test Street",
				IsDefault:     true,
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().UpdateAddress(
					mock.Anything,
					&user.UpdateAddressRequest{
						AddressId:     1,
						ReceiverName:  "John Doe",
						ReceiverPhone: "1234567890",
						Province:      "TestProvince",
						City:          "TestCity",
						District:      "TestDistrict",
						DetailAddress: "123 Test Street",
						IsDefault:     true,
					},
				).Return(&user.UpdateAddressResponse{}, nil)
			},
			wantErr: nil,
		},
		{
			name: "address not found",
			req: &types.AddressReq{
				AddressId:     999,
				ReceiverName:  "John Doe",
				ReceiverPhone: "1234567890",
				Province:      "TestProvince",
				City:          "TestCity",
				District:      "TestDistrict",
				DetailAddress: "123 Test Street",
				IsDefault:     false,
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().UpdateAddress(
					mock.Anything,
					&user.UpdateAddressRequest{
						AddressId:     999,
						ReceiverName:  "John Doe",
						ReceiverPhone: "1234567890",
						Province:      "TestProvince",
						City:          "TestCity",
						District:      "TestDistrict",
						DetailAddress: "123 Test Street",
						IsDefault:     false,
					},
				).Return(nil, errors.New("address not found"))
			},
			wantErr: errors.New("address not found"),
		},
		{
			name: "rpc error",
			req: &types.AddressReq{
				AddressId:     1,
				ReceiverName:  "John Doe",
				ReceiverPhone: "1234567890",
				Province:      "TestProvince",
				City:          "TestCity",
				District:      "TestDistrict",
				DetailAddress: "123 Test Street",
				IsDefault:     false,
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().UpdateAddress(
					mock.Anything,
					&user.UpdateAddressRequest{
						AddressId:     1,
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

			logic := NewUpdateAddressLogic(context.Background(), svcCtx)
			err := logic.UpdateAddress(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
