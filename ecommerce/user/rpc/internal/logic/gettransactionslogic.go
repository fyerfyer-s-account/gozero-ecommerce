package logic

import (
	"context"
	"math"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTransactionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTransactionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTransactionsLogic {
	return &GetTransactionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTransactionsLogic) GetTransactions(in *user.GetTransactionsRequest) (*user.GetTransactionsResponse, error) {
	// todo: add your logic here and delete this line
	// 1. Validate page parameters
	if in.Page < 1 || in.PageSize < 1 {
		return nil, zeroerr.ErrInvalidTransactionParams
	}
	if in.PageSize > 100 {
		in.PageSize = 100
	}

	// 2. Check if user exists
	_, err := l.svcCtx.UsersModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, zeroerr.ErrUserNotFound
		}
		return nil, err
	}

	// 3. Get transactions with filter
	var transactions []*model.WalletTransactions
	if in.Type > 0 {
		transactions, err = l.svcCtx.WalletTransactionsModel.FindByType(l.ctx, uint64(in.UserId), int64(in.Type))
		if err != nil {
			logx.Errorf("get transactions by type error: %v", err)
			return nil, zeroerr.ErrGetTransactionsFailed
		}
	} else {
		transactions, err = l.svcCtx.WalletTransactionsModel.FindByUserId(l.ctx, uint64(in.UserId), int(in.Page), int(in.PageSize))
		if err != nil {
			logx.Errorf("get transactions error: %v", err)
			return nil, zeroerr.ErrGetTransactionsFailed
		}
	}

	// 4. Convert model to proto
	resp := make([]*user.Transaction, 0, len(transactions))
	for _, t := range transactions {
		resp = append(resp, &user.Transaction{
			Id:        int64(t.Id),
			UserId:    int64(t.UserId),
			Amount:    t.Amount,
			Type:      t.Type,
			Status:    t.Status,
			Remark:    t.Remark.String,
			CreatedAt: t.CreatedAt.Unix(),
		})
	}

	// 5. Calculate total pages
	total := int64(len(transactions))
	totalPages := int32(math.Ceil(float64(total) / float64(in.PageSize)))

	return &user.GetTransactionsResponse{
		Transactions: resp,
		Total:        total,
		Page:         in.Page,
		TotalPages:   totalPages,
	}, nil
}
