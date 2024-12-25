package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSkusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListSkusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSkusLogic {
	return &ListSkusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListSkusLogic) ListSkus(in *product.ListSkusRequest) (*product.ListSkusResponse, error) {
	// todo: add your logic here and delete this line

	return &product.ListSkusResponse{}, nil
}
