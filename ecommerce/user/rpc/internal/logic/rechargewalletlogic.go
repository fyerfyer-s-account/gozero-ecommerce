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

type RechargeWalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRechargeWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RechargeWalletLogic {
	return &RechargeWalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RechargeWalletLogic) RechargeWallet(in *user.RechargeWalletRequest) (*user.RechargeWalletResponse, error) {
	// 1. Validate input
	if err := l.validateRechargeInput(in); err != nil {
		return nil, err
	}

	var (
		orderId string
		balance float64
	)

	// 2. Execute recharge transaction
	err := l.svcCtx.WalletAccountsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// Generate order ID
		orderId = fmt.Sprintf("R%d%d", in.UserId, time.Now().UnixNano())

		// Create transaction record
		transResult, err := l.svcCtx.WalletTransactionsModel.WithSession(session).Insert(ctx, &model.WalletTransactions{
			UserId:  uint64(in.UserId),
			OrderId: orderId,
			Amount:  in.Amount,
			Type:    1, // Recharge
			Status:  3, // Processing
			Remark:  sql.NullString{String: "钱包充值-" + in.Channel, Valid: true},
		})
		if err != nil {
			return err
		}

		// Update wallet balance
		err = l.svcCtx.WalletAccountsModel.WithSession(session).UpdateBalance(ctx, uint64(in.UserId), in.Amount)
		if err != nil {
			return err
		}

		// Get updated wallet
		wallet, err := l.svcCtx.WalletAccountsModel.WithSession(session).FindOneByUserId(ctx, uint64(in.UserId))
		if err != nil {
			return err
		}
		balance = wallet.Balance

		// Update transaction status
		transId, _ := transResult.LastInsertId()
		err = l.svcCtx.WalletTransactionsModel.WithSession(session).Update(ctx, &model.WalletTransactions{
			Id:     uint64(transId),
			Status: 1, // Success
		})
		return err
	})

	if err != nil {
		logx.Errorf("recharge wallet error: %v", err)
		return nil, zeroerr.ErrRechargeWalletFailed
	}

	return &user.RechargeWalletResponse{
		OrderId: orderId,
		Balance: balance,
	}, nil
}

func (l *RechargeWalletLogic) validateRechargeInput(in *user.RechargeWalletRequest) error {
	// Check user exists
	_, err := l.svcCtx.UsersModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		if err == model.ErrNotFound {
			return zeroerr.ErrUserNotFound
		}
		return err
	}

	// Validate amount
	if in.Amount <= 0 {
		return zeroerr.ErrInvalidAmount
	}

	// Check wallet exists and status
	wallet, err := l.svcCtx.WalletAccountsModel.FindOneByUserId(l.ctx, uint64(in.UserId))
	if err != nil && err != model.ErrNotFound {
		return err
	}
	if wallet != nil && wallet.Status != 1 {
		return zeroerr.ErrWalletDisabled
	}

	return nil
}
