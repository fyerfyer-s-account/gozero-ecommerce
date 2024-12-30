package logic

import (
	"context"
	"flag"
	"strings"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestRechargeWalletLogic_RechargeWallet(t *testing.T) {
	configFile := flag.String("f", "../../etc/user.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test user
	testUser := &model.Users{
		Username: "testrecharge",
		Password: "testpass123",
		Status:   1,
	}
	result, err := ctx.UsersModel.Insert(context.Background(), testUser)
	assert.NoError(t, err)
	userId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create wallet account
	wallet := &model.WalletAccounts{
		UserId:  uint64(userId),
		Balance: 0,
		Status:  1,
	}
	_, err = ctx.WalletAccountsModel.Insert(context.Background(), wallet)
	assert.NoError(t, err)

	// Cleanup function
	defer func() {
		ctx.UsersModel.Delete(context.Background(), uint64(userId))
		ctx.WalletAccountsModel.DeleteByUserId(context.Background(), uint64(userId))
		transactions, _ := ctx.WalletTransactionsModel.FindByUserId(context.Background(), uint64(userId), 1, 100)
		for _, trans := range transactions {
			ctx.WalletTransactionsModel.Delete(context.Background(), trans.Id)
		}
	}()

	tests := []struct {
		name    string
		req     *user.RechargeWalletRequest
		wantErr error
		check   func(*testing.T, *user.RechargeWalletResponse)
	}{
		{
			name: "Successful recharge",
			req: &user.RechargeWalletRequest{
				UserId:  userId,
				Amount:  100.00,
				Channel: "alipay",
			},
			wantErr: nil,
			check: func(t *testing.T, resp *user.RechargeWalletResponse) {
				assert.NotEmpty(t, resp.OrderId)
				assert.True(t, strings.HasPrefix(resp.OrderId, "R"))
				assert.Equal(t, 100.00, resp.Balance)

				// Verify wallet balance
				wallet, err := ctx.WalletAccountsModel.FindOneByUserId(context.Background(), uint64(userId))
				assert.NoError(t, err)
				assert.Equal(t, 100.00, wallet.Balance)

				// Verify transaction record
				trans, err := ctx.WalletTransactionsModel.FindByUserId(context.Background(), uint64(userId), 1, 1)
				assert.NoError(t, err)
				assert.Len(t, trans, 1)
				assert.Equal(t, 100.00, trans[0].Amount)
				assert.Equal(t, int64(1), trans[0].Status)
			},
		},
		{
			name: "Invalid amount",
			req: &user.RechargeWalletRequest{
				UserId:  userId,
				Amount:  -100.00,
				Channel: "alipay",
			},
			wantErr: zeroerr.ErrInvalidAmount,
		},
		{
			name: "Non-existent user",
			req: &user.RechargeWalletRequest{
				UserId:  99999,
				Amount:  100.00,
				Channel: "alipay",
			},
			wantErr: zeroerr.ErrUserNotFound,
		},
		{
			name: "Zero amount",
			req: &user.RechargeWalletRequest{
				UserId:  userId,
				Amount:  0,
				Channel: "alipay",
			},
			wantErr: zeroerr.ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewRechargeWalletLogic(context.Background(), ctx)
			resp, err := l.RechargeWallet(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				if tt.check != nil {
					tt.check(t, resp)
				}
			}
		})
	}
}
