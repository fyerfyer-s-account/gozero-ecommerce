package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTemplateLogic {
	return &UpdateTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateTemplateLogic) UpdateTemplate(in *message.UpdateTemplateRequest) (*message.UpdateTemplateResponse, error) {
	// todo: add your logic here and delete this line

	return &message.UpdateTemplateResponse{}, nil
}
