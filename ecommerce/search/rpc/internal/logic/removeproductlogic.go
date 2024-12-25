package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveProductLogic {
	return &RemoveProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveProductLogic) RemoveProduct(in *search.RemoveProductRequest) (*search.RemoveProductResponse, error) {
	// todo: add your logic here and delete this line

	return &search.RemoveProductResponse{}, nil
}
