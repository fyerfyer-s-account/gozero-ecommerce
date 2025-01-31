package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/util"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
    // Check input validation
    if len(in.UserIds) == 0 {
        return &message.BatchSendMessageResponse{}, nil
    }

    if in.Type <= 0 || in.SendChannel <= 0 {
        return nil, zeroerr.ErrInvalidMessageType
    }

    resp := &message.BatchSendMessageResponse{
        MessageIds: make([]int64, 0),
        Errors:     make([]*message.BatchSendError, 0),
    }

    // Create messages in batch
    messages := make([]*model.Messages, 0, len(in.UserIds))
    for _, userId := range in.UserIds {
        msg := &model.Messages{
            UserId:      uint64(userId),
            Title:       in.Title,
            Content:     in.Content,
            Type:        int64(in.Type),
            SendChannel: int64(in.SendChannel),
            IsRead:      0,
        }
        if in.ExtraData != "" {
            msg.ExtraData = sql.NullString{
                String: in.ExtraData,
                Valid:  true,
            }
        }
        messages = append(messages, msg)
    }

    // Insert messages in transaction
    err := l.svcCtx.MessagesModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
        return l.svcCtx.MessagesModel.BatchInsert(ctx, messages)
    })
    if err != nil {
        logx.Errorf("Failed to batch insert messages: %v", err)
        return nil, zeroerr.ErrMessageCreateFailed
    }

    // Create batch message event
    batchId := uuid.New().String()
    batchEvent := &types.MessageBatchEvent{
        MessageEvent: types.MessageEvent{
            Type:      types.MessageBatchCreated,
            UserID:    in.UserIds[0], // Use first user as reference
            Timestamp: time.Now(),
        },
        BatchID:   batchId,
        Total:     int32(len(messages)),
        Completed: 0,
        Failed:    0,
        Status:    "pending",
    }

    // Publish batch created event
    if err := l.svcCtx.Producer.PublishBatchEvent(l.ctx, batchEvent); err != nil {
        logx.Errorf("Failed to publish batch created event: %v", err)
    }

    // Send individual messages
    for _, msg := range messages {
        resp.MessageIds = append(resp.MessageIds, int64(msg.Id))

        // Create individual message sent event
        messageEvent := &types.MessageEventSentEvent{
            MessageEvent: types.MessageEvent{
                Type:      types.MessageEventSent,
                UserID:    int64(msg.UserId),
                Timestamp: time.Now(),
            },
            MessageID: fmt.Sprintf("%d", msg.Id),
            Channel:   util.GetChannelString(int32(msg.SendChannel)),
            Content:   msg.Content,
            Recipient: fmt.Sprintf("%d", msg.UserId),
            Status:    "pending",
            Variables: map[string]string{
                "title":     msg.Title,
                "type":      fmt.Sprintf("%d", msg.Type),
                "extraData": msg.ExtraData.String,
                "batchId":   batchId,
            },
        }

        if err := l.svcCtx.Producer.PublishMessageEvent(l.ctx, messageEvent); err != nil {
            logx.Errorf("Failed to publish message %d: %v", msg.Id, err)
            resp.Errors = append(resp.Errors, &message.BatchSendError{
                UserId: int64(msg.UserId),
                Error:  zeroerr.ErrSendMessageFailed.Error(),
            })
        }
    }

    // If all messages failed, return batch send failure
    if len(resp.Errors) == len(messages) {
        return nil, zeroerr.ErrBatchSendFailed
    }

    return resp, nil
}