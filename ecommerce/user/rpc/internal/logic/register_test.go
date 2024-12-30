package logic

import (
	"context"
	"flag"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestRegisterLogic_Register(t *testing.T) {
	configFile := flag.String("f", "../../etc/user.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Setup test cases
	tests := []struct {
		name    string
		req     *user.RegisterRequest
		wantErr error
		cleanup bool // indicate if cleanup is needed
	}{
		{
			name: "Valid registration",
			req: &user.RegisterRequest{
				Username: "testuser",
				Password: "password123",
				Phone:    "13800138000",
				Email:    "test@example.com",
			},
			wantErr: nil,
			cleanup: true,
		},
		{
			name: "Empty username",
			req: &user.RegisterRequest{
				Username: "",
				Password: "password123",
			},
			wantErr: zeroerr.ErrUsernameTooShort,
			cleanup: false,
		},
		{
			name: "Invalid phone",
			req: &user.RegisterRequest{
				Username: "testuser2",
				Password: "password123",
				Phone:    "invalid",
			},
			wantErr: zeroerr.ErrInvalidPhone,
			cleanup: false,
		},
		{
			name: "Invalid email",
			req: &user.RegisterRequest{
				Username: "testuser3",
				Password: "password123",
				Email:    "invalid-email",
			},
			wantErr: zeroerr.ErrInvalidEmail,
			cleanup: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create logic instance
			l := NewRegisterLogic(context.Background(), ctx)

			// Execute test
			resp, err := l.Register(tt.req)

			// Verify results
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Greater(t, resp.UserId, int64(0))

				// Cleanup if test case requires it
				if tt.cleanup {
					// Delete the created user
					err := ctx.UsersModel.Delete(context.Background(), uint64(resp.UserId))
					assert.NoError(t, err, "Failed to cleanup test user")

					// Delete wallet account
					err = ctx.WalletAccountsModel.DeleteByUserId(context.Background(),
						uint64(resp.UserId))
					assert.NoError(t, err, "Failed to cleanup wallet account")

					// Delete wallet transactions
					err = ctx.WalletTransactionsModel.Delete(context.Background(),
						uint64(resp.UserId))
					assert.NoError(t, err, "Failed to cleanup wallet transactions")
				}
			}
		})
	}
}
