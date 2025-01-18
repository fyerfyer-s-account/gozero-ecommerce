package logic

import (
    "context"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/types"
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

    // Publish message read event
    readTime := time.Now()
    metadata := types.Metadata{
        Source: "message-service",
        UserID: in.UserId,
    }
    // 安全获取trace_id
    if traceID, ok := l.ctx.Value("trace_id").(string); ok {
        metadata.TraceID = traceID
    }
    
    err = l.svcCtx.Producer.PublishMessageRead(l.ctx, &types.MessageReadData{
        MessageID: in.MessageId,
        UserID:    in.UserId,
        ReadTime:  readTime,
    }, metadata)
    if err != nil {
        logx.Errorf("Failed to publish message read event: %v", err)
        // Don't return error as the message is already marked as read
    }

    return &message.ReadMessageResponse{
        Success: true,
    }, nil
}