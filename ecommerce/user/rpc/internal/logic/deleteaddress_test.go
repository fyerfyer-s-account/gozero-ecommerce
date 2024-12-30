package logic

import (
	"context"
	"flag"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestDeleteAddressLogic_DeleteAddress(t *testing.T) {
	// Load config
	configFile := flag.String("f", "../../etc/user.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test user
	testUser := &model.Users{
		Username: "testdeleteaddr",
		Password: "testpass123",
		Status:   1,
	}
	result, err := ctx.UsersModel.Insert(context.Background(), testUser)
	assert.NoError(t, err)
	userId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create test addresses
	defaultAddr := &model.UserAddresses{
		UserId:        uint64(userId),
		ReceiverName:  "Default User",
		ReceiverPhone: "13800138000",
		Province:      "Test Province",
		City:          "Test City",
		District:      "Test District",
		DetailAddress: "Default Address",
		IsDefault:     1,
	}
	defaultResult, err := ctx.UserAddressesModel.Insert(context.Background(), defaultAddr)
	assert.NoError(t, err)
	defaultAddrId, err := defaultResult.LastInsertId()
	assert.NoError(t, err)

	normalAddr := &model.UserAddresses{
		UserId:        uint64(userId),
		ReceiverName:  "Normal User",
		ReceiverPhone: "13800138001",
		Province:      "Test Province",
		City:          "Test City",
		District:      "Test District",
		DetailAddress: "Normal Address",
		IsDefault:     0,
	}
	normalResult, err := ctx.UserAddressesModel.Insert(context.Background(), normalAddr)
	assert.NoError(t, err)
	normalAddrId, err := normalResult.LastInsertId()
	assert.NoError(t, err)

	// Cleanup
	defer func() {
		err := ctx.UsersModel.Delete(context.Background(), uint64(userId))
		assert.NoError(t, err)
		err = ctx.UserAddressesModel.Delete(context.Background(), uint64(defaultAddrId))
		assert.NoError(t, err)
		err = ctx.UserAddressesModel.Delete(context.Background(), uint64(normalAddrId))
		assert.NoError(t, err)
	}()

	tests := []struct {
		name    string
		req     *user.DeleteAddressRequest
		wantErr error
	}{
		{
			name: "Delete normal address",
			req: &user.DeleteAddressRequest{
				UserId:    userId,
				AddressId: normalAddrId,
			},
			wantErr: nil,
		},
		{
			name: "Try delete default address",
			req: &user.DeleteAddressRequest{
				UserId:    userId,
				AddressId: defaultAddrId,
			},
			wantErr: zeroerr.ErrDefaultAddressNotDeletable,
		},
		{
			name: "Delete non-existent address",
			req: &user.DeleteAddressRequest{
				UserId:    userId,
				AddressId: 99999,
			},
			wantErr: zeroerr.ErrAddressNotFound,
		},
		{
			name: "Delete with invalid user",
			req: &user.DeleteAddressRequest{
				UserId:    99999,
				AddressId: normalAddrId,
			},
			wantErr: zeroerr.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewDeleteAddressLogic(context.Background(), ctx)
			resp, err := l.DeleteAddress(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)

				// Verify address was actually deleted
				_, err := ctx.UserAddressesModel.FindOne(context.Background(), uint64(tt.req.AddressId))
				assert.Equal(t, model.ErrNotFound, err)
			}
		})
	}
}
