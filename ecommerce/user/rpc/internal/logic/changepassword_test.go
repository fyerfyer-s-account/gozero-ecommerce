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

func TestChangePasswordLogic_ChangePassword(t *testing.T) {
	// Load config
	configFile := flag.String("f", "../../etc/user.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test user
	initialPassword := "oldpassword123"
	hashedPassword := cryptx.HashPassword(initialPassword, c.Salt)
	testUser := &model.Users{
		Username: "testpassworduser",
		Password: hashedPassword,
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
		req     *user.ChangePasswordRequest
		wantErr error
	}{
		{
			name: "Invalid old password",
			req: &user.ChangePasswordRequest{
				UserId:      userId,
				OldPassword: "wrongpassword",
				NewPassword: "newpassword123",
			},
			wantErr: zeroerr.ErrOldPasswordIncorrect,
		},
		{
			name: "Same old and new password",
			req: &user.ChangePasswordRequest{
				UserId:      userId,
				OldPassword: initialPassword,
				NewPassword: initialPassword,
			},
			wantErr: zeroerr.ErrSamePassword,
		},
		{
			name: "Password too short",
			req: &user.ChangePasswordRequest{
				UserId:      userId,
				OldPassword: initialPassword,
				NewPassword: "short",
			},
			wantErr: zeroerr.ErrPasswordTooWeak,
		},
		{
			name: "Non-existent user",
			req: &user.ChangePasswordRequest{
				UserId:      99999,
				OldPassword: initialPassword,
				NewPassword: "newpassword123",
			},
			wantErr: zeroerr.ErrUserNotFound,
		},
		{
			name: "Valid password change",
			req: &user.ChangePasswordRequest{
				UserId:      userId,
				OldPassword: initialPassword,
				NewPassword: "newpassword123",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewChangePasswordLogic(context.Background(), ctx)
			resp, err := l.ChangePassword(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)

				// Verify password was actually changed
				updatedUser, err := ctx.UsersModel.FindOne(context.Background(), uint64(userId))
				assert.NoError(t, err)
				assert.Equal(t, cryptx.HashPassword(tt.req.NewPassword, c.Salt), updatedUser.Password)
			}
		})
	}
}
