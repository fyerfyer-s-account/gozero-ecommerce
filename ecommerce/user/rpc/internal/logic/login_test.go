package logic

import (
	"context"
	"flag"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/cryptx"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestLoginLogic_Login(t *testing.T) {
	// Load config
	configFile := flag.String("f", "../../etc/user.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test user for login tests
	testUser := &model.Users{
		Username: "testloginuser",
		Password: cryptx.HashPassword("testpassword123", c.Salt),
		Status:   1,
	}
	result, err := ctx.UsersModel.Insert(context.Background(), testUser)
	assert.NoError(t, err)
	userId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Cleanup function
	defer func() {
		err := ctx.UsersModel.Delete(context.Background(), uint64(userId))
		assert.NoError(t, err)
	}()

	tests := []struct {
		name    string
		req     *user.LoginRequest
		wantErr error
	}{
		{
			name: "Valid login",
			req: &user.LoginRequest{
				Username: "testloginuser",
				Password: "testpassword123",
			},
			wantErr: nil,
		},
		{
			name: "Empty username",
			req: &user.LoginRequest{
				Username: "",
				Password: "testpassword123",
			},
			wantErr: zeroerr.ErrInvalidUsername,
		},
		{
			name: "Invalid password",
			req: &user.LoginRequest{
				Username: "testloginuser",
				Password: "wrongpassword",
			},
			wantErr: zeroerr.ErrInvalidPassword,
		},
		{
			name: "Non-existent account",
			req: &user.LoginRequest{
				Username: "nonexistentuser",
				Password: "testpassword123",
			},
			wantErr: zeroerr.ErrInvalidAccount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoginLogic(context.Background(), ctx)
			resp, err := l.Login(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.AccessToken)
				assert.NotEmpty(t, resp.RefreshToken)
				assert.Greater(t, resp.ExpiresIn, int64(0))
			}
		})
	}
}
