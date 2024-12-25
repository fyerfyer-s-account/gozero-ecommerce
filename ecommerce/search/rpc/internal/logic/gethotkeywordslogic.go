package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHotKeywordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetHotKeywordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHotKeywordsLogic {
	return &GetHotKeywordsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetHotKeywordsLogic) GetHotKeywords(in *search.GetHotKeywordsRequest) (*search.GetHotKeywordsResponse, error) {
	// todo: add your logic here and delete this line

	return &search.GetHotKeywordsResponse{}, nil
}
