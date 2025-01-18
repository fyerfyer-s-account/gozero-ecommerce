package logic

import (
    "context"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
    "github.com/zeromicro/go-zero/core/logx"
)

type GetMessageLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewGetMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageLogic {
    return &GetMessageLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *GetMessageLogic) GetMessage(in *message.GetMessageRequest) (*message.GetMessageResponse, error) {
    // Get message from database
    msg, err := l.svcCtx.MessagesModel.FindOne(l.ctx, uint64(in.MessageId))
    if err != nil {
        return nil, zeroerr.ErrMessageNotFound
    }

    // Convert to proto message
    return &message.GetMessageResponse{
        Message: &message.Message{
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
        },
    }, nil
}