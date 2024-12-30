package logic

import (
	"context"
	"database/sql"
	"flag"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestGetTransactionsLogic_GetTransactions(t *testing.T) {
	// Load config
	configFile := flag.String("f", "../../etc/user.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test user
	testUser := &model.Users{
		Username: "testtransactions",
		Password: "testpass123",
		Status:   1,
	}
	result, err := ctx.UsersModel.Insert(context.Background(), testUser)
	assert.NoError(t, err)
	userId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create test transactions
	transactions := []*model.WalletTransactions{
		{
			UserId:    uint64(userId),
			Amount:    100.00,
			OrderId:   "abcdefg",
			Type:      1, // deposit
			Status:    1,
			Remark:    sql.NullString{String: "Test deposit", Valid: true},
			CreatedAt: time.Now(),
		},
		{
			UserId:    uint64(userId),
			Amount:    -50.00,
			OrderId:   "hijklmn",
			Type:      2, // withdrawal
			Status:    1,
			Remark:    sql.NullString{String: "Test withdrawal", Valid: true},
			CreatedAt: time.Now(),
		},
	}

	var transactionIds []int64
	for _, tx := range transactions {
		res, err := ctx.WalletTransactionsModel.Insert(context.Background(), tx)
		assert.NoError(t, err)
		id, err := res.LastInsertId()
		assert.NoError(t, err)
		transactionIds = append(transactionIds, id)
	}

	// Cleanup
	defer func() {
		err := ctx.UsersModel.Delete(context.Background(), uint64(userId))
		assert.NoError(t, err)
		for _, id := range transactionIds {
			err := ctx.WalletTransactionsModel.Delete(context.Background(), uint64(id))
			assert.NoError(t, err)
		}
	}()

	tests := []struct {
		name    string
		req     *user.GetTransactionsRequest
		wantErr error
		check   func(*testing.T, *user.GetTransactionsResponse)
	}{
		{
			name: "Get all transactions",
			req: &user.GetTransactionsRequest{
				UserId:   userId,
				Page:     1,
				PageSize: 10,
			},
			wantErr: nil,
			check: func(t *testing.T, resp *user.GetTransactionsResponse) {
				assert.Equal(t, int64(2), resp.Total)
				assert.Equal(t, int32(1), resp.TotalPages)
				assert.Len(t, resp.Transactions, 2)
			},
		},
		{
			name: "Get transactions by type",
			req: &user.GetTransactionsRequest{
				UserId:   userId,
				Page:     1,
				PageSize: 10,
				Type:     1,
			},
			wantErr: nil,
			check: func(t *testing.T, resp *user.GetTransactionsResponse) {
				assert.Equal(t, int64(1), resp.Total)
				assert.Len(t, resp.Transactions, 1)
				assert.Equal(t, int64(1), resp.Transactions[0].Type)
			},
		},
		{
			name: "Invalid page parameters",
			req: &user.GetTransactionsRequest{
				UserId:   userId,
				Page:     0,
				PageSize: 0,
			},
			wantErr: zeroerr.ErrInvalidTransactionParams,
			check:   nil,
		},
		{
			name: "Non-existent user",
			req: &user.GetTransactionsRequest{
				UserId:   99999,
				Page:     1,
				PageSize: 10,
			},
			wantErr: zeroerr.ErrUserNotFound,
			check:   nil,
		},
		{
			name: "Page size exceeding limit",
			req: &user.GetTransactionsRequest{
				UserId:   userId,
				Page:     1,
				PageSize: 200,
			},
			wantErr: nil,
			check: func(t *testing.T, resp *user.GetTransactionsResponse) {
				assert.Equal(t, int64(2), resp.Total)
				assert.Equal(t, int32(1), resp.Page)
				assert.Len(t, resp.Transactions, 2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewGetTransactionsLogic(context.Background(), ctx)
			resp, err := l.GetTransactions(tt.req)

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
