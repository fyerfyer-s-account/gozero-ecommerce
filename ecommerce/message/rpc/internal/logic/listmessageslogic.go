package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMessagesLogic {
	return &ListMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMessagesLogic) ListMessages(in *message.ListMessagesRequest) (*message.ListMessagesResponse, error) {
	// todo: add your logic here and delete this line

	return &message.ListMessagesResponse{}, nil
}
