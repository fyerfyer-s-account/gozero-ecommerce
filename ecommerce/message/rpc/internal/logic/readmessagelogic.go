package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReadMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadMessageLogic {
	return &ReadMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReadMessageLogic) ReadMessage(in *message.ReadMessageRequest) (*message.ReadMessageResponse, error) {
	// todo: add your logic here and delete this line

	return &message.ReadMessageResponse{}, nil
}
