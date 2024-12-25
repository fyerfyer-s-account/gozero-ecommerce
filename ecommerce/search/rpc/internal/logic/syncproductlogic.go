package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncProductLogic {
	return &SyncProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 索引管理
func (l *SyncProductLogic) SyncProduct(in *search.SyncProductRequest) (*search.SyncProductResponse, error) {
	// todo: add your logic here and delete this line

	return &search.SyncProductResponse{}, nil
}
