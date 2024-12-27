package user

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"

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

	return
}
