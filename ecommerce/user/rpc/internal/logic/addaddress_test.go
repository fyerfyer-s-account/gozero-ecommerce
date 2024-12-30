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

func TestAddAddressLogic_AddAddress(t *testing.T) {
	// Load config
	configFile := flag.String("f", "../../etc/user.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test user
	testUser := &model.Users{
		Username: "testaddressuser",
		Password: "testpass123",
		Status:   1,
	}
	result, err := ctx.UsersModel.Insert(context.Background(), testUser)
	assert.NoError(t, err)
	userId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Cleanup function
	defer func() {
		// Delete test user
		err := ctx.UsersModel.Delete(context.Background(), uint64(userId))
		assert.NoError(t, err)

		// Delete test addresses
		addresses, _ := ctx.UserAddressesModel.FindByUserId(context.Background(), uint64(userId))
		for _, addr := range addresses {
			err := ctx.UserAddressesModel.Delete(context.Background(), addr.Id)
			assert.NoError(t, err)
		}
	}()

	tests := []struct {
		name    string
		req     *user.AddAddressRequest
		wantErr error
	}{
		{
			name: "Valid address addition",
			req: &user.AddAddressRequest{
				UserId:        userId,
				ReceiverName:  "Test User",
				ReceiverPhone: "13800138000",
				Province:      "Test Province",
				City:          "Test City",
				District:      "Test District",
				DetailAddress: "Test Detail Address",
				IsDefault:     true,
			},
			wantErr: nil,
		},
		{
			name: "Invalid user ID",
			req: &user.AddAddressRequest{
				UserId:        99999,
				ReceiverName:  "Test User",
				ReceiverPhone: "13800138000",
				Province:      "Test Province",
				City:          "Test City",
				District:      "Test District",
				DetailAddress: "Test Detail Address",
			},
			wantErr: zeroerr.ErrUserNotFound,
		},
		{
			name: "Invalid address input",
			req: &user.AddAddressRequest{
				UserId:        userId,
				ReceiverName:  "",
				ReceiverPhone: "",
				Province:      "",
				City:          "",
				District:      "",
				DetailAddress: "",
			},
			wantErr: zeroerr.ErrInvalidAddress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewAddAddressLogic(context.Background(), ctx)
			resp, err := l.AddAddress(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Greater(t, resp.AddressId, int64(0))

				// Verify address was actually created
				addr, err := ctx.UserAddressesModel.FindOne(context.Background(), uint64(resp.AddressId))
				assert.NoError(t, err)
				assert.Equal(t, tt.req.ReceiverName, addr.ReceiverName)
				assert.Equal(t, tt.req.ReceiverPhone, addr.ReceiverPhone)
			}
		})
	}
}
