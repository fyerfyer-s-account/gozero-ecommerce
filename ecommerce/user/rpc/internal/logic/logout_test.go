package logic

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/cryptx"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/jwtx"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestLogoutLogic_Logout(t *testing.T) {
	configFile := flag.String("f", "../../etc/user.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test user
	testUser := &model.Users{
		Username: "testlogoutuser",
		Password: cryptx.HashPassword("testpassword123", c.Salt),
		Online:   1, // Initially online
		Status:   1, // Account enabled
	}
	result, err := ctx.UsersModel.Insert(context.Background(), testUser)
	assert.NoError(t, err)
	userId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Generate test token
	token, err := jwtx.GetToken(c.JwtAuth.AccessSecret, time.Now().Unix(), c.JwtAuth.AccessExpire, userId, jwtx.RoleUser)
	assert.NoError(t, err)

	defer func() {
		err := ctx.UsersModel.Delete(context.Background(), uint64(userId))
		assert.NoError(t, err)
	}()

	tests := []struct {
		name    string
		req     *user.LogoutRequest
		wantErr error
	}{
		{
			name: "Valid logout",
			req: &user.LogoutRequest{
				AccessToken: token,
			},
			wantErr: nil,
		},
		{
			name: "Invalid token",
			req: &user.LogoutRequest{
				AccessToken: "invalid_token",
			},
			wantErr: zeroerr.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLogoutLogic(context.Background(), ctx)
			resp, err := l.Logout(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)

				// Verify user status is updated
				user, err := ctx.UsersModel.FindOne(context.Background(), uint64(userId))
				assert.NoError(t, err)
				assert.Equal(t, int32(0), user.Status)
				assert.Equal(t, int32(0), user.Online)
			}
		})
	}
}

