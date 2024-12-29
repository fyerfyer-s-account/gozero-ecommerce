package user

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWalletTransactionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWalletTransactionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWalletTransactionsLogic {
	return &GetWalletTransactionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWalletTransactionsLogic) GetWalletTransactions(req *types.TransactionListReq) (resp *types.TransactionListResp, err error) {
	// todo: add your logic here and delete this line
	// Get userId from JWT context
	userId := l.ctx.Value("userId").(int64)

	// Call RPC
	res, err := l.svcCtx.UserRpc.GetTransactions(l.ctx, &user.GetTransactionsRequest{
		UserId:   userId,
		Page:     req.Page,
		PageSize: req.PageSize,
		Type:     req.Type,
	})

	if err != nil {
		logx.Errorf("get wallet transactions error: %v", err)
		return nil, err
	}

	// Convert transactions
	transactions := make([]types.Transaction, 0, len(res.Transactions))
	for _, t := range res.Transactions {
		transactions = append(transactions, types.Transaction{
			Id:        t.Id,
			UserId:    t.UserId,
			OrderId:   t.OrderId,
			Amount:    t.Amount,
			Type:      t.Type,
			Status:    t.Status,
			Remark:    t.Remark,
			CreatedAt: t.CreatedAt,
		})
	}

	return &types.TransactionListResp{
		List:       transactions,
		Total:      res.Total,
		Page:       res.Page,
		TotalPages: res.TotalPages,
	}, nil
}
