package logic

import (
	"context"
	"flag"
	"log"
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
	log.Println("Loaded configuration file successfully.")

	ctx := svc.NewServiceContext(c)
	log.Println("Initialized service context.")

	// Create test user
	testUser := &model.Users{
		Username: "testrecharge",
		Password: "testpass123",
		Status:   1,
	}
	result, err := ctx.UsersModel.Insert(context.Background(), testUser)
	assert.NoError(t, err, "Failed to insert test user.")
	userId, err := result.LastInsertId()
	assert.NoError(t, err, "Failed to get user ID.")
	log.Printf("Test user created with ID: %d", userId)

	// Create wallet account
	wallet := &model.WalletAccounts{
		UserId:  uint64(userId),
		Balance: 0,
		Status:  1,
	}
	_, err = ctx.WalletAccountsModel.Insert(context.Background(), wallet)
	assert.NoError(t, err, "Failed to insert wallet account.")
	log.Printf("Wallet account created for user ID: %d", userId)

	// Cleanup function
	defer func() {
		log.Println("Cleaning up test data.")
		ctx.UsersModel.Delete(context.Background(), uint64(userId))
		ctx.WalletAccountsModel.DeleteByUserId(context.Background(), uint64(userId))
		transactions, _ := ctx.WalletTransactionsModel.FindByUserId(context.Background(), uint64(userId), 1, 100)
		for _, trans := range transactions {
			ctx.WalletTransactionsModel.Delete(context.Background(), trans.Id)
		}
		log.Println("Cleanup complete.")
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
				log.Println("Validating successful recharge response.")
				assert.NotEmpty(t, resp.OrderId, "Order ID is empty.")
				assert.True(t, strings.HasPrefix(resp.OrderId, "R"), "Order ID does not have the correct prefix.")
				assert.Equal(t, 100.00, resp.Balance, "Balance mismatch.")

				// Verify wallet balance
				wallet, err := ctx.WalletAccountsModel.FindOneByUserId(context.Background(), uint64(userId))
				assert.NoError(t, err, "Failed to fetch wallet account.")
				assert.Equal(t, 100.00, wallet.Balance, "Wallet balance mismatch.")

				// Verify transaction record
				trans, err := ctx.WalletTransactionsModel.FindByUserId(context.Background(), uint64(userId), 1, 1)
				assert.NoError(t, err, "Failed to fetch transaction records.")
				assert.Len(t, trans, 1, "Transaction record count mismatch.")
				assert.Equal(t, 100.00, trans[0].Amount, "Transaction amount mismatch.")
				assert.Equal(t, int64(1), trans[0].Status, "Transaction status mismatch.")
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
		log.Printf("Running test case: %s", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			l := NewRechargeWalletLogic(context.Background(), ctx)
			resp, err := l.RechargeWallet(tt.req)

			if tt.wantErr != nil {
				log.Printf("Expected error: %v", tt.wantErr)
				assert.Error(t, err, "Expected an error but got none.")
				assert.Equal(t, tt.wantErr, err, "Error mismatch.")
				assert.Nil(t, resp, "Expected response to be nil.")
			} else {
				assert.NoError(t, err, "Unexpected error occurred.")
				assert.NotNil(t, resp, "Expected a response but got nil.")
				if tt.check != nil {
					tt.check(t, resp)
				}
			}
		})
	}
}
