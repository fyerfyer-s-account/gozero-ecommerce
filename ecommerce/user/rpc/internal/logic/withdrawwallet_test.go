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

func TestWithdrawWalletLogic_WithdrawWallet(t *testing.T) {
	// Load config
	configFile := flag.String("f", "../../etc/user.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test user
	testUser := &model.Users{
		Username: "testwithdraw",
		Status:   1,
	}
	result, err := ctx.UsersModel.Insert(context.Background(), testUser)
	assert.NoError(t, err)
	userId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create wallet with initial balance
	wallet := &model.WalletAccounts{
		UserId:  uint64(userId),
		Balance: 1000.00,
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
		req     *user.WithdrawWalletRequest
		wantErr error
	}{
		{
			name: "Successful withdrawal",
			req: &user.WithdrawWalletRequest{
				UserId:   userId,
				Amount:   100.00,
				BankCard: "6222000000000000",
			},
			wantErr: nil,
		},
		{
			name: "Invalid amount",
			req: &user.WithdrawWalletRequest{
				UserId:   userId,
				Amount:   0,
				BankCard: "6222000000000000",
			},
			wantErr: zeroerr.ErrInvalidAmount,
		},
		{
			name: "Insufficient balance",
			req: &user.WithdrawWalletRequest{
				UserId:   userId,
				Amount:   2000.00,
				BankCard: "6222000000000000",
			},
			wantErr: zeroerr.ErrInsufficientBalance,
		},
		{
			name: "Invalid bank card",
			req: &user.WithdrawWalletRequest{
				UserId:   userId,
				Amount:   100.00,
				BankCard: "",
			},
			wantErr: zeroerr.ErrInvalidBankCard,
		},
		{
			name: "Non-existent user",
			req: &user.WithdrawWalletRequest{
				UserId:   99999,
				Amount:   100.00,
				BankCard: "6222000000000000",
			},
			wantErr: zeroerr.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewWithdrawWalletLogic(context.Background(), ctx)
			resp, err := l.WithdrawWallet(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.OrderId)

				// Verify wallet balance update
				updatedWallet, err := ctx.WalletAccountsModel.FindOneByUserId(context.Background(), uint64(tt.req.UserId))
				assert.NoError(t, err)
				assert.Equal(t, 900.00, updatedWallet.Balance)

				// Verify transaction record
				trans, err := ctx.WalletTransactionsModel.FindOneByOrderId(context.Background(), resp.OrderId)
				assert.NoError(t, err)
				assert.Equal(t, uint64(tt.req.UserId), trans.UserId)
				assert.Equal(t, tt.req.Amount, trans.Amount)
				assert.Equal(t, int64(2), trans.Type)
				assert.Equal(t, int64(1), trans.Status)
			}
		})
	}
}
