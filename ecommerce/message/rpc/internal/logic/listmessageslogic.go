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
    // Get messages with filters
    messages, err := l.svcCtx.MessagesModel.FindByUserId(
        l.ctx,
        uint64(in.UserId),
        int64(in.Type),
        in.UnreadOnly,
        int(in.Page),
        int(in.PageSize),
    )
    if err != nil {
        logx.Errorf("Failed to get messages: %v", err)
        return nil, err
    }

    // Get total count
    total, err := l.svcCtx.MessagesModel.CountByUserId(
        l.ctx,
        uint64(in.UserId),
        int64(in.Type),
        in.UnreadOnly,
    )
    if err != nil {
        logx.Errorf("Failed to get total count: %v", err)
        return nil, err
    }

    // Convert to proto messages
    protoMessages := make([]*message.Message, 0, len(messages))
    for _, msg := range messages {
        protoMessages = append(protoMessages, &message.Message{
            Id:          int64(msg.Id),
            UserId:      int64(msg.UserId),
            Title:       msg.Title,
            Content:     msg.Content,
            Type:        int32(msg.Type),
            SendChannel: int32(msg.SendChannel),
            ExtraData:   msg.ExtraData.String,
            IsRead:      msg.IsRead == 1,
            ReadTime:    msg.ReadTime.Time.Unix(),
            CreatedAt:   msg.CreatedAt.Unix(),
        })
    }

    return &message.ListMessagesResponse{
        Messages: protoMessages,
        Total:    total,
    }, nil
}