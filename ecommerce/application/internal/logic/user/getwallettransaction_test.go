package user

import (
	"context"
	"errors"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetWalletTransactionsLogic_GetWalletTransactions(t *testing.T) {
	tests := []struct {
		name    string
		userId  int64
		req     *types.TransactionListReq
		mock    func(mockUser *User)
		want    *types.TransactionListResp
		wantErr error
	}{
		{
			name:   "successful transactions retrieval",
			userId: 12345,
			req: &types.TransactionListReq{
				Page:     1,
				PageSize: 10,
				Type:     1,
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().GetTransactions(
					mock.Anything,
					&user.GetTransactionsRequest{
						UserId:   12345,
						Page:     1,
						PageSize: 10,
						Type:     1,
					},
				).Return(&user.GetTransactionsResponse{
					Transactions: []*user.Transaction{
						{
							Id:        1,
							UserId:    12345,
							OrderId:   "ORDER123",
							Amount:    100.50,
							Type:      1,
							Status:    1,
							Remark:    "Payment",
							CreatedAt: 1634567890,
						},
					},
					Total:      1,
					Page:       1,
					TotalPages: 1,
				}, nil)
			},
			want: &types.TransactionListResp{
				List: []types.Transaction{
					{
						Id:        1,
						UserId:    12345,
						OrderId:   "ORDER123",
						Amount:    100.50,
						Type:      1,
						Status:    1,
						Remark:    "Payment",
						CreatedAt: 1634567890,
					},
				},
				Total:      1,
				Page:       1,
				TotalPages: 1,
			},
			wantErr: nil,
		},
		{
			name:   "empty transaction list",
			userId: 12345,
			req: &types.TransactionListReq{
				Page:     1,
				PageSize: 10,
				Type:     1,
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().GetTransactions(
					mock.Anything,
					&user.GetTransactionsRequest{
						UserId:   12345,
						Page:     1,
						PageSize: 10,
						Type:     1,
					},
				).Return(&user.GetTransactionsResponse{
					Transactions: []*user.Transaction{},
					Total:        0,
					Page:         1,
					TotalPages:   0,
				}, nil)
			},
			want: &types.TransactionListResp{
				List:       []types.Transaction{},
				Total:      0,
				Page:       1,
				TotalPages: 0,
			},
			wantErr: nil,
		},
		{
			name:   "rpc error",
			userId: 12345,
			req: &types.TransactionListReq{
				Page:     1,
				PageSize: 10,
				Type:     1,
			},
			mock: func(mockUser *User) {
				mockUser.EXPECT().GetTransactions(
					mock.Anything,
					&user.GetTransactionsRequest{
						UserId:   12345,
						Page:     1,
						PageSize: 10,
						Type:     1,
					},
				).Return(nil, errors.New("rpc error"))
			},
			want:    nil,
			wantErr: errors.New("rpc error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := NewUser(t)
			tt.mock(mockUser)

			svcCtx := &svc.ServiceContext{
				UserRpc: mockUser,
			}

			ctx := context.WithValue(context.Background(), "userId", tt.userId)
			logic := NewGetWalletTransactionsLogic(ctx, svcCtx)

			got, err := logic.GetWalletTransactions(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
