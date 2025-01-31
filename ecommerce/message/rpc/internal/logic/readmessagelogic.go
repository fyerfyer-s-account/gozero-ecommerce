package logic

import (
    "context"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
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
    // Verify message exists
    msg, err := l.svcCtx.MessagesModel.FindOne(l.ctx, uint64(in.MessageId))
    if err != nil {
        return nil, zeroerr.ErrMessageNotFound
    }

    // Verify message belongs to user
    if msg.UserId != uint64(in.UserId) {
        return nil, zeroerr.ErrMessageNotFound
    }

    // Update read status
    err = l.svcCtx.MessagesModel.UpdateReadStatus(l.ctx, uint64(in.MessageId), uint64(in.UserId))
    if err != nil {
        return nil, zeroerr.ErrMessageUpdateFailed
    }

    return &message.ReadMessageResponse{
        Success: true,
    }, nil
}