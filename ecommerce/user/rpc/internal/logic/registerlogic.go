package logic

import (
	"context"
	"database/sql"
	"regexp"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/cryptx"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户注册
func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	// 1. Input validation
	if len(in.Username) == 0 || len(in.Password) < l.svcCtx.Config.MinPasswordLength {
		return nil, zeroerr.ErrInvalidUsername
	}

	if in.Phone != "" && !validatePhone(in.Phone) {
		return nil, zeroerr.ErrInvalidPhone
	}

	if in.Email != "" && !validateEmail(in.Email) {
		return nil, zeroerr.ErrInvalidEmail
	}

	// 2. Check existing user
	exist, err := l.checkExistingUser(l.ctx, in)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, zeroerr.ErrDuplicateUsername
	}

	// 3. Create user transaction
	var userId uint64
	err = l.svcCtx.UsersModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// Hash password
		hashedPassword := cryptx.HashPassword(in.Password, l.svcCtx.Config.Salt)

		// Insert user
		result, err := l.svcCtx.UsersModel.WithSession(session).Insert(ctx, &model.Users{
			Username: in.Username,
			Password: hashedPassword,
			Phone:    sql.NullString{String: in.Phone, Valid: len(in.Phone) > 0},
			Email:    sql.NullString{String: in.Email, Valid: len(in.Email) > 0},
			Status:   1, // Active
		})
		if err != nil {
			return err
		}

		insertId, err := result.LastInsertId()
		if err != nil {
			return err
		}
		userId = uint64(insertId)

		// Create wallet account
		_, err = l.svcCtx.WalletAccountsModel.WithSession(session).Insert(ctx, &model.WalletAccounts{
			UserId:  userId,
			Balance: 0,
		})
		if err != nil {
			return err
		}

		// Initialize user points if configured
		if l.svcCtx.Config.InitialPoints > 0 {
			_, err = l.svcCtx.WalletTransactionsModel.WithSession(session).Insert(ctx, &model.WalletTransactions{
				UserId: userId,
				Amount: float64(l.svcCtx.Config.InitialPoints),
				Type:   1, // Points
				Status: 1, // Success
				Remark: sql.NullString{
					String: "新用户注册奖励",
					Valid:  true,
				},
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, zeroerr.ErrRegisterFailed
	}

	return &user.RegisterResponse{
		UserId: int64(userId),
	}, nil
}

func (l *RegisterLogic) checkExistingUser(ctx context.Context, in *user.RegisterRequest) (bool, error) {
	if _, err := l.svcCtx.UsersModel.FindOneByUsername(ctx, in.Username); err == nil {
		return true, nil
	}

	if in.Phone != "" {
		if _, err := l.svcCtx.UsersModel.FindOneByPhone(ctx, sql.NullString{String: in.Phone, Valid: true}); err == nil {
			return true, nil
		}
	}

	if in.Email != "" {
		if _, err := l.svcCtx.UsersModel.FindOneByEmail(ctx, sql.NullString{String: in.Email, Valid: true}); err == nil {
			return true, nil
		}
	}

	return false, nil
}

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func validatePhone(phone string) bool {
	// start with 13x,14x,15x,16x,17x,18x,19x
	pattern := `^1[3-9]\d{9}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(phone)
}
