package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTemplatesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTemplatesLogic {
	return &ListTemplatesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListTemplatesLogic) ListTemplates(in *message.ListTemplatesRequest) (*message.ListTemplatesResponse, error) {
	// todo: add your logic here and delete this line

	return &message.ListTemplatesResponse{}, nil
}
