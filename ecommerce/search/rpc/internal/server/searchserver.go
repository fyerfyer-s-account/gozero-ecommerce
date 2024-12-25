// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3
// Source: search.proto

package server

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/internal/logic"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/search"
)

type SearchServer struct {
	svcCtx *svc.ServiceContext
	search.UnimplementedSearchServer
}

func NewSearchServer(svcCtx *svc.ServiceContext) *SearchServer {
	return &SearchServer{
		svcCtx: svcCtx,
	}
}

// 商品搜索
func (s *SearchServer) SearchProducts(ctx context.Context, in *search.SearchProductsRequest) (*search.SearchProductsResponse, error) {
	l := logic.NewSearchProductsLogic(ctx, s.svcCtx)
	return l.SearchProducts(in)
}

func (s *SearchServer) GetHotKeywords(ctx context.Context, in *search.GetHotKeywordsRequest) (*search.GetHotKeywordsResponse, error) {
	l := logic.NewGetHotKeywordsLogic(ctx, s.svcCtx)
	return l.GetHotKeywords(in)
}

func (s *SearchServer) GetSearchSuggestions(ctx context.Context, in *search.GetSearchSuggestionsRequest) (*search.GetSearchSuggestionsResponse, error) {
	l := logic.NewGetSearchSuggestionsLogic(ctx, s.svcCtx)
	return l.GetSearchSuggestions(in)
}

func (s *SearchServer) SaveSearchHistory(ctx context.Context, in *search.SaveSearchHistoryRequest) (*search.SaveSearchHistoryResponse, error) {
	l := logic.NewSaveSearchHistoryLogic(ctx, s.svcCtx)
	return l.SaveSearchHistory(in)
}

func (s *SearchServer) GetSearchHistory(ctx context.Context, in *search.GetSearchHistoryRequest) (*search.GetSearchHistoryResponse, error) {
	l := logic.NewGetSearchHistoryLogic(ctx, s.svcCtx)
	return l.GetSearchHistory(in)
}

func (s *SearchServer) DeleteSearchHistory(ctx context.Context, in *search.DeleteSearchHistoryRequest) (*search.DeleteSearchHistoryResponse, error) {
	l := logic.NewDeleteSearchHistoryLogic(ctx, s.svcCtx)
	return l.DeleteSearchHistory(in)
}

// 索引管理
func (s *SearchServer) SyncProduct(ctx context.Context, in *search.SyncProductRequest) (*search.SyncProductResponse, error) {
	l := logic.NewSyncProductLogic(ctx, s.svcCtx)
	return l.SyncProduct(in)
}

func (s *SearchServer) RemoveProduct(ctx context.Context, in *search.RemoveProductRequest) (*search.RemoveProductResponse, error) {
	l := logic.NewRemoveProductLogic(ctx, s.svcCtx)
	return l.RemoveProduct(in)
}
