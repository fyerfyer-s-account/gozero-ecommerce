package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveSearchHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveSearchHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveSearchHistoryLogic {
	return &SaveSearchHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SaveSearchHistoryLogic) SaveSearchHistory(in *search.SaveSearchHistoryRequest) (*search.SaveSearchHistoryResponse, error) {
	// todo: add your logic here and delete this line

	return &search.SaveSearchHistoryResponse{}, nil
}
