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

func TestListAddressesLogic_ListAddresses(t *testing.T) {
	tests := []struct {
		name    string
		userId  int64
		mock    func(mockUser *User)
		want    []types.Address
		wantErr error
	}{
		{
			name:   "successful addresses retrieval",
			userId: 12345,
			mock: func(mockUser *User) {
				mockUser.EXPECT().GetUserAddresses(
					mock.Anything,
					&user.GetUserAddressesRequest{
						UserId: 12345,
					},
				).Return(&user.GetUserAddressesResponse{
					Addresses: []*user.Address{
						{
							Id:            1,
							ReceiverName:  "John Doe",
							ReceiverPhone: "1234567890",
							Province:      "TestProvince",
							City:          "TestCity",
							District:      "TestDistrict",
							DetailAddress: "123 Test Street",
							IsDefault:     true,
						},
						{
							Id:            2,
							ReceiverName:  "Jane Doe",
							ReceiverPhone: "0987654321",
							Province:      "TestProvince2",
							City:          "TestCity2",
							District:      "TestDistrict2",
							DetailAddress: "456 Test Avenue",
							IsDefault:     false,
						},
					},
				}, nil)
			},
			want: []types.Address{
				{
					Id:            1,
					ReceiverName:  "John Doe",
					ReceiverPhone: "1234567890",
					Province:      "TestProvince",
					City:          "TestCity",
					District:      "TestDistrict",
					DetailAddress: "123 Test Street",
					IsDefault:     true,
				},
				{
					Id:            2,
					ReceiverName:  "Jane Doe",
					ReceiverPhone: "0987654321",
					Province:      "TestProvince2",
					City:          "TestCity2",
					District:      "TestDistrict2",
					DetailAddress: "456 Test Avenue",
					IsDefault:     false,
				},
			},
			wantErr: nil,
		},
		{
			name:   "empty address list",
			userId: 12345,
			mock: func(mockUser *User) {
				mockUser.EXPECT().GetUserAddresses(
					mock.Anything,
					&user.GetUserAddressesRequest{
						UserId: 12345,
					},
				).Return(&user.GetUserAddressesResponse{
					Addresses: []*user.Address{},
				}, nil)
			},
			want:    []types.Address{},
			wantErr: nil,
		},
		{
			name:   "rpc error",
			userId: 12345,
			mock: func(mockUser *User) {
				mockUser.EXPECT().GetUserAddresses(
					mock.Anything,
					&user.GetUserAddressesRequest{
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
			mockUser := NewUser(t)
			tt.mock(mockUser)

			svcCtx := &svc.ServiceContext{
				UserRpc: mockUser,
			}

			ctx := context.WithValue(context.Background(), "userId", tt.userId)
			logic := NewListAddressesLogic(ctx, svcCtx)

			got, err := logic.ListAddresses()

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
