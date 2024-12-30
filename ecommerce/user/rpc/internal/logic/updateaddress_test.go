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

func TestUpdateAddressLogic_UpdateAddress(t *testing.T) {
	// Load config
	configFile := flag.String("f", "../../etc/user.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test user
	testUser := &model.Users{
		Username: "testuser",
		Status:   1,
	}
	result, err := ctx.UsersModel.Insert(context.Background(), testUser)
	assert.NoError(t, err)
	userId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create test address
	testAddress := &model.UserAddresses{
		UserId:        uint64(userId),
		ReceiverName:  "John Doe",
		ReceiverPhone: "13800138000",
		Province:      "Test Province",
		City:          "Test City",
		District:      "Test District",
		DetailAddress: "Test Detail Address",
		IsDefault:     0,
	}
	addrResult, err := ctx.UserAddressesModel.Insert(context.Background(), testAddress)
	assert.NoError(t, err)
	addressId, err := addrResult.LastInsertId()
	assert.NoError(t, err)

	// Cleanup
	defer func() {
		err := ctx.UserAddressesModel.Delete(context.Background(), uint64(addressId))
		assert.NoError(t, err)
		err = ctx.UsersModel.Delete(context.Background(), uint64(userId))
		assert.NoError(t, err)
	}()

	tests := []struct {
		name    string
		req     *user.UpdateAddressRequest
		wantErr error
	}{
		{
			name: "Successful update",
			req: &user.UpdateAddressRequest{
				AddressId:     addressId,
				ReceiverName:  "Jane Doe",
				ReceiverPhone: "13900139000",
				Province:      "New Province",
				City:          "New City",
				District:      "New District",
				DetailAddress: "New Detail Address",
				IsDefault:     true,
			},
			wantErr: nil,
		},
		{
			name: "Invalid address ID",
			req: &user.UpdateAddressRequest{
				AddressId:     99999,
				ReceiverName:  "Jane Doe",
				ReceiverPhone: "13900139000",
				Province:      "New Province",
				City:          "New City",
				District:      "New District",
				DetailAddress: "New Detail Address",
			},
			wantErr: zeroerr.ErrAddressNotFound,
		},
		{
			name: "Invalid input - empty receiver name",
			req: &user.UpdateAddressRequest{
				AddressId:     addressId,
				ReceiverName:  "",
				ReceiverPhone: "13900139000",
				Province:      "New Province",
				City:          "New City",
				District:      "New District",
				DetailAddress: "New Detail Address",
			},
			wantErr: zeroerr.ErrInvalidAddress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewUpdateAddressLogic(context.Background(), ctx)
			resp, err := l.UpdateAddress(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)

				// Verify address was actually updated
				updatedAddr, err := ctx.UserAddressesModel.FindOne(context.Background(), uint64(tt.req.AddressId))
				assert.NoError(t, err)
				assert.Equal(t, tt.req.ReceiverName, updatedAddr.ReceiverName)
				assert.Equal(t, tt.req.ReceiverPhone, updatedAddr.ReceiverPhone)
				assert.Equal(t, tt.req.Province, updatedAddr.Province)
				assert.Equal(t, tt.req.City, updatedAddr.City)
				assert.Equal(t, tt.req.District, updatedAddr.District)
				assert.Equal(t, tt.req.DetailAddress, updatedAddr.DetailAddress)
				assert.Equal(t, boolToInt64(tt.req.IsDefault), updatedAddr.IsDefault)
			}
		})
	}
}
