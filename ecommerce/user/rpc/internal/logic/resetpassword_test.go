package logic

import (
	"context"
	"database/sql"
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

func TestResetPasswordLogic_ResetPassword(t *testing.T) {
	// Load config
	configFile := flag.String("f", "../../etc/user.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test user
	testPhone := "13800138000"
	testUser := &model.Users{
		Username: "testresetpwd",
		Password: cryptx.HashPassword("oldpassword", c.Salt),
		Phone:    sql.NullString{String: testPhone, Valid: true},
		Status:   1,
	}
	result, err := ctx.UsersModel.Insert(context.Background(), testUser)
	assert.NoError(t, err)
	userId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Set verify code in Redis
	verifyCode := "123456"
	codeKey := fmt.Sprintf("reset:code:%s", testPhone)
	err = ctx.BizRedis.Set(codeKey, verifyCode)
	assert.NoError(t, err)

	// Cleanup
	defer func() {
		err := ctx.UsersModel.Delete(context.Background(), uint64(userId))
		assert.NoError(t, err)
		ctx.BizRedis.Del(codeKey)
	}()

	tests := []struct {
		name    string
		req     *user.ResetPasswordRequest
		wantErr error
		setup   func()
	}{
		{
			name: "Successful reset",
			req: &user.ResetPasswordRequest{
				Phone:       testPhone,
				VerifyCode:  verifyCode,
				NewPassword: "newpassword123",
			},
			wantErr: nil,
			setup: func() {
				ctx.BizRedis.Set(codeKey, verifyCode)
			},
		},
		{
			name: "Invalid phone",
			req: &user.ResetPasswordRequest{
				Phone:       "13900139000",
				VerifyCode:  verifyCode,
				NewPassword: "newpassword123",
			},
			wantErr: zeroerr.ErrPhoneNotFound,
			setup: func() {
				ctx.BizRedis.Set(codeKey, verifyCode)
			},
		},
		{
			name: "Invalid verify code",
			req: &user.ResetPasswordRequest{
				Phone:       testPhone,
				VerifyCode:  "wrong_code",
				NewPassword: "newpassword123",
			},
			wantErr: zeroerr.ErrInvalidVerifyCode,
			setup: func() {
				ctx.BizRedis.Set(codeKey, verifyCode)
			},
		},
		{
			name: "Weak password",
			req: &user.ResetPasswordRequest{
				Phone:       testPhone,
				VerifyCode:  verifyCode,
				NewPassword: "123",
			},
			wantErr: zeroerr.ErrPasswordTooWeak,
			setup: func() {
				ctx.BizRedis.Set(codeKey, verifyCode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset Redis state before each test
			ctx.BizRedis.Del(codeKey)
			// Setup test state
			tt.setup()

			l := NewResetPasswordLogic(context.Background(), ctx)
			resp, err := l.ResetPassword(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)

				// Verify password was actually changed
				user, err := ctx.UsersModel.FindOneByPhone(context.Background(),
					sql.NullString{String: tt.req.Phone, Valid: true})
				assert.NoError(t, err)
				assert.NotEqual(t, cryptx.HashPassword("oldpassword", c.Salt),
					user.Password)
			}
		})
	}
}
