package logic

import (
    "context"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
    "github.com/zeromicro/go-zero/core/logx"
)

type DeleteMessageLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewDeleteMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessageLogic {
    return &DeleteMessageLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *DeleteMessageLogic) DeleteMessage(in *message.DeleteMessageRequest) (*message.DeleteMessageResponse, error) {
    // Verify message exists and belongs to user
    msg, err := l.svcCtx.MessagesModel.FindOne(l.ctx, uint64(in.MessageId))
    if err != nil {
        return nil, zeroerr.ErrMessageNotFound
    }

    if msg.UserId != uint64(in.UserId) {
        return nil, zeroerr.ErrMessageNotFound
    }

    // Delete message
    err = l.svcCtx.MessagesModel.DeleteByUserMessage(l.ctx, uint64(in.MessageId), uint64(in.UserId))
    if err != nil {
        logx.Errorf("Failed to delete message: %v", err)
        return nil, err
    }

    return &message.DeleteMessageResponse{
        Success: true,
    }, nil
}