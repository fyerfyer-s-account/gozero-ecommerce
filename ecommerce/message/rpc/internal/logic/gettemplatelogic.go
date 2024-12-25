package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTemplateLogic {
	return &GetTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTemplateLogic) GetTemplate(in *message.GetTemplateRequest) (*message.GetTemplateResponse, error) {
	// todo: add your logic here and delete this line

	return &message.GetTemplateResponse{}, nil
}
