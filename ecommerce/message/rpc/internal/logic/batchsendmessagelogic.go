package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchSendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchSendMessageLogic {
	return &BatchSendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchSendMessageLogic) BatchSendMessage(in *message.BatchSendMessageRequest) (*message.BatchSendMessageResponse, error) {
	// todo: add your logic here and delete this line

	return &message.BatchSendMessageResponse{}, nil
}
