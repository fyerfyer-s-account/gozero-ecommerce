package logic

import (
	"context"
	"flag"
	"fmt"
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
	configFile := flag.String("f", "../../etc/user.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test user for login tests
	testUser := &model.Users{
		Username: "testloginuser",
		Password: cryptx.HashPassword("testpassword123", c.Salt),
		Online:   0, // Initially offline
		Status:   1, // Account enabled
	}
	result, err := ctx.UsersModel.Insert(context.Background(), testUser)
	assert.NoError(t, err)
	userId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create already logged in user
	loggedInUser := &model.Users{
		Username: "loggedinuser",
		Password: cryptx.HashPassword("testpassword123", c.Salt),
		Status:   1, // Already logged in
	}
	resultLoggedIn, err := ctx.UsersModel.Insert(context.Background(), loggedInUser)
	assert.NoError(t, err)
	loggedInUserId, err := resultLoggedIn.LastInsertId()
	assert.NoError(t, err)

	// Create already online user
	onlineUser := &model.Users{
		Username: "onlineuser",
		Password: cryptx.HashPassword("testpassword123", c.Salt),
		Online:   1, // Already online
		Status:   1,
	}
	resultOnline, err := ctx.UsersModel.Insert(context.Background(), onlineUser)
	assert.NoError(t, err)
	onlineUserId, err := resultOnline.LastInsertId()
	assert.NoError(t, err)

	// Cleanup function
	defer func() {
		err := ctx.UsersModel.Delete(context.Background(), uint64(userId))
		assert.NoError(t, err)
		err = ctx.UsersModel.Delete(context.Background(), uint64(loggedInUserId))
		assert.NoError(t, err)
		err = ctx.UsersModel.Delete(context.Background(), uint64(onlineUserId))
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
			name: "Already logged in user",
			req: &user.LoginRequest{
				Username: "loggedinuser",
				Password: "testpassword123",
			},
			wantErr: zeroerr.ErrAlreadyLoggedIn,
		},
		{
			name: "Already online user",
			req: &user.LoginRequest{
				Username: "onlineuser",
				Password: "testpassword123",
			},
			wantErr: zeroerr.ErrAlreadyLoggedIn,
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

				// Verify user status is updated to logged in
				user, err := ctx.UsersModel.FindOne(context.Background(), uint64(userId))
				assert.NoError(t, err)
				assert.Equal(t, int32(1), user.Status)
				assert.Equal(t, int32(1), user.Online)

				// Verify old token is invalidated
				oldTokenKey := fmt.Sprintf("%s%d", c.JwtAuth.RefreshRedis.KeyPrefix, userId)
				_, err = ctx.RefreshRedis.Get(oldTokenKey)
				assert.Error(t, err) // Should be redis.Nil error
			}
		})
	}
}
