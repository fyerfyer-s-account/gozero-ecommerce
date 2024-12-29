package logic

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"time"

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
	// Log input
	log.Printf("Register called with Username: %s, Phone: %s, Email: %s", in.Username, in.Phone, in.Email)

	// 1. Input validation
	if len(in.Username) == 0 {
		log.Println("Validation failed: Username too short")
		return nil, zeroerr.ErrUsernameTooShort
	}
	if len(in.Password) < l.svcCtx.Config.MinPasswordLength {
		log.Printf("Validation failed: Password too short, length: %d", len(in.Password))
		return nil, zeroerr.ErrPasswordTooShort
	}
	if in.Phone != "" && !validatePhone(in.Phone) {
		log.Printf("Validation failed: Invalid phone format: %s", in.Phone)
		return nil, zeroerr.ErrInvalidPhone
	}
	if in.Email != "" && !validateEmail(in.Email) {
		log.Printf("Validation failed: Invalid email format: %s", in.Email)
		return nil, zeroerr.ErrInvalidEmail
	}

	// 2. Check existing user
	log.Println("Checking if user already exists")
	exist, err := l.checkExistingUser(l.ctx, in)
	if err != nil {
		log.Printf("Error checking existing user: %v", err)
		return nil, zeroerr.ErrRegisterFailed
	}
	if exist {
		log.Println("User already exists")
		return nil, zeroerr.ErrUserExists
	}

	// 3. Create user transaction
	var userId uint64
	err = l.svcCtx.UsersModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		log.Println("Starting user creation transaction")

		// Hash password
		hashedPassword := cryptx.HashPassword(in.Password, l.svcCtx.Config.Salt)
		log.Println("Password hashed successfully")

		// Insert user
		log.Println("Inserting user into database")
		result, err := l.svcCtx.UsersModel.WithSession(session).Insert(ctx, &model.Users{
			Username: in.Username,
			Password: hashedPassword,
			Phone:    sql.NullString{String: in.Phone, Valid: len(in.Phone) > 0},
			Email:    sql.NullString{String: in.Email, Valid: len(in.Email) > 0},
			Status:   1, // Active
		})
		if err != nil {
			log.Printf("Error inserting user: %v", err)
			return err
		}

		insertId, err := result.LastInsertId()
		if err != nil {
			log.Printf("Error getting last insert ID: %v", err)
			return err
		}
		userId = uint64(insertId)
		log.Printf("User created with ID: %d", userId)

		// Create wallet account
		log.Println("Creating wallet account for user")
		_, err = l.svcCtx.WalletAccountsModel.WithSession(session).Insert(ctx, &model.WalletAccounts{
			UserId:  userId,
			Balance: 0,
		})
		if err != nil {
			log.Printf("Error creating wallet account: %v", err)
			return zeroerr.ErrInitWalletFailed
		}

		// Initialize user points if configured
		if l.svcCtx.Config.InitialPoints > 0 {
			log.Printf("Initializing points: %d", l.svcCtx.Config.InitialPoints)
			orderId := fmt.Sprintf("W%d%d", userId, time.Now().UnixNano())
			_, err = l.svcCtx.WalletTransactionsModel.WithSession(session).Insert(ctx, &model.WalletTransactions{
				UserId:  userId,
				OrderId: orderId,
				Amount:  float64(l.svcCtx.Config.InitialPoints),
				Type:    1, // Points
				Status:  1, // Success
				Remark: sql.NullString{
					String: "新用户注册奖励",
					Valid:  true,
				},
			})
			if err != nil {
				log.Printf("Error initializing points: %v", err)
				return zeroerr.ErrInitPointsFailed
			}
		}

		log.Println("Transaction completed successfully")
		return nil
	})

	if err != nil {
		log.Printf("Error in transaction: %v", err)
		switch err {
		case zeroerr.ErrInitWalletFailed:
			return nil, zeroerr.ErrInitWalletFailed
		case zeroerr.ErrInitPointsFailed:
			return nil, zeroerr.ErrInitPointsFailed
		default:
			return nil, zeroerr.ErrRegisterFailed
		}
	}

	log.Printf("User registration successful, UserID: %d", userId)
	return &user.RegisterResponse{
		UserId: int64(userId),
	}, nil
}

func (l *RegisterLogic) checkExistingUser(ctx context.Context, in *user.RegisterRequest) (bool, error) {
	if _, err := l.svcCtx.UsersModel.FindOneByUsername(ctx, in.Username); err == nil {
		return true, zeroerr.ErrDuplicateUsername
	}

	if in.Phone != "" {
		if _, err := l.svcCtx.UsersModel.FindOneByPhone(ctx, sql.NullString{String: in.Phone, Valid: true}); err == nil {
			return true, zeroerr.ErrDuplicatePhone
		}
	}

	if in.Email != "" {
		if _, err := l.svcCtx.UsersModel.FindOneByEmail(ctx, sql.NullString{String: in.Email, Valid: true}); err == nil {
			return true, zeroerr.ErrDuplicateEmail
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
