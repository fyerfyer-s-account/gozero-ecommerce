package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendTemplateMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendTemplateMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendTemplateMessageLogic {
	return &SendTemplateMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendTemplateMessageLogic) SendTemplateMessage(in *message.SendTemplateMessageRequest) (*message.SendTemplateMessageResponse, error) {
	// todo: add your logic here and delete this line

	return &message.SendTemplateMessageResponse{}, nil
}
