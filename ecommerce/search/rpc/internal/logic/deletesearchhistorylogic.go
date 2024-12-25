package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSearchHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSearchHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSearchHistoryLogic {
	return &DeleteSearchHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteSearchHistoryLogic) DeleteSearchHistory(in *search.DeleteSearchHistoryRequest) (*search.DeleteSearchHistoryResponse, error) {
	// todo: add your logic here and delete this line

	return &search.DeleteSearchHistoryResponse{}, nil
}
