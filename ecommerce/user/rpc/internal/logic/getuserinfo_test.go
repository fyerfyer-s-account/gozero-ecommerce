package logic

import (
	"context"
	"database/sql"
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

func TestGetUserInfoLogic_GetUserInfo(t *testing.T) {
	// Load config
	configFile := flag.String("f", "../../etc/user.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test user
	testUser := &model.Users{
		Username: "testuser",
		Password: "testpass",
		Nickname: sql.NullString{String: "Test User", Valid: true},
		Phone:    sql.NullString{String: "1234567890", Valid: true},
		Email:    sql.NullString{String: "test@example.com", Valid: true},
		Status:   1,
	}
	result, err := ctx.UsersModel.Insert(context.Background(), testUser)
	assert.NoError(t, err)
	userId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create wallet for test user
	wallet := &model.WalletAccounts{
		UserId:  uint64(userId),
		Balance: 100.00,
		Status:  1,
	}
	_, err = ctx.WalletAccountsModel.Insert(context.Background(), wallet)
	assert.NoError(t, err)

	// Cleanup
	defer func() {
		err := ctx.WalletAccountsModel.DeleteByUserId(context.Background(), uint64(userId))
		assert.NoError(t, err)
		err = ctx.UsersModel.Delete(context.Background(), uint64(userId))
		assert.NoError(t, err)
	}()

	tests := []struct {
		name    string
		req     *user.GetUserInfoRequest
		wantErr error
	}{
		{
			name: "Valid user info",
			req: &user.GetUserInfoRequest{
				UserId: userId,
			},
			wantErr: nil,
		},
		{
			name: "Non-existent user",
			req: &user.GetUserInfoRequest{
				UserId: 99999,
			},
			wantErr: zeroerr.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewGetUserInfoLogic(context.Background(), ctx)
			resp, err := l.GetUserInfo(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.req.UserId, resp.UserInfo.UserId)
				assert.Equal(t, "testuser", resp.UserInfo.Username)
				assert.Equal(t, "Test User", resp.UserInfo.Nickname)
				assert.Equal(t, "1234567890", resp.UserInfo.Phone)
				assert.Equal(t, "test@example.com", resp.UserInfo.Email)
				assert.Equal(t, float64(100), resp.UserInfo.WalletBalance)
			}
		})
	}
}
