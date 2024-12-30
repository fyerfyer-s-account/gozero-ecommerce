package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type WithdrawWalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func (l *WithdrawWalletLogic) WithdrawWallet(in *user.WithdrawWalletRequest) (*user.WithdrawWalletResponse, error) {
	// 1. Validate input
	if err := l.validateWithdrawInput(in); err != nil {
		return nil, err
	}

	var (
		orderId string
		balance float64
	)

	// 2. Execute withdraw transaction
	err := l.svcCtx.WalletAccountsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// Generate order ID
		orderId = fmt.Sprintf("W%d%d", in.UserId, time.Now().UnixNano())

		// Create transaction record
		transResult, err := l.svcCtx.WalletTransactionsModel.WithSession(session).Insert(ctx, &model.WalletTransactions{
			UserId:  uint64(in.UserId),
			OrderId: orderId,
			Amount:  in.Amount,
			Type:    2, // Withdraw
			Status:  3, // Processing
			Remark:  sql.NullString{String: "钱包提现-" + in.BankCard, Valid: true},
		})
		if err != nil {
			return err
		}

		// Freeze amount for withdrawal
		err = l.svcCtx.WalletAccountsModel.WithSession(session).FreezeAmount(ctx, uint64(in.UserId), in.Amount)
		if err != nil {
			return err
		}

		// Get wallet balance
		wallet, err := l.svcCtx.WalletAccountsModel.WithSession(session).FindOneByUserId(ctx, uint64(in.UserId))
		if err != nil {
			return err
		}
		balance = wallet.Balance

		// Update transaction status
		transId, _ := transResult.LastInsertId()
		_, err = l.svcCtx.WalletTransactionsModel.WithSession(session).UpdateState(ctx, 1, uint64(transId))
		return err
	})

	if err != nil {
		logx.Errorf("withdraw wallet error: %v", err)
		return nil, zeroerr.ErrWithdrawFailed
	}

	return &user.WithdrawWalletResponse{
		OrderId: orderId,
		Balance: balance,
	}, nil
}

func NewWithdrawWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawWalletLogic {
	return &WithdrawWalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WithdrawWalletLogic) validateWithdrawInput(in *user.WithdrawWalletRequest) error {
	// Check if user exists and get wallet
	wallet, err := l.svcCtx.WalletAccountsModel.FindOneByUserId(l.ctx, uint64(in.UserId))
	if err != nil {
		if err == model.ErrNotFound {
			return zeroerr.ErrUserNotFound
		}
		return err
	}

	// Check wallet status
	if wallet.Status != 1 {
		return zeroerr.ErrWalletDisabled
	}

	// Validate amount
	if in.Amount <= 0 {
		return zeroerr.ErrInvalidAmount
	}

	// Check if has sufficient balance
	if wallet.Balance < in.Amount {
		return zeroerr.ErrInsufficientBalance
	}

	// Validate bank card
	if len(in.BankCard) == 0 {
		return zeroerr.ErrInvalidBankCard
	}

	return nil
}
