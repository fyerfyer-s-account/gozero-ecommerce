package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTemplateLogic {
	return &CreateTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 消息模板
func (l *CreateTemplateLogic) CreateTemplate(in *message.CreateTemplateRequest) (*message.CreateTemplateResponse, error) {
	// todo: add your logic here and delete this line

	return &message.CreateTemplateResponse{}, nil
}
