package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSearchSuggestionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSearchSuggestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSearchSuggestionsLogic {
	return &GetSearchSuggestionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSearchSuggestionsLogic) GetSearchSuggestions(in *search.GetSearchSuggestionsRequest) (*search.GetSearchSuggestionsResponse, error) {
	// todo: add your logic here and delete this line

	return &search.GetSearchSuggestionsResponse{}, nil
}
