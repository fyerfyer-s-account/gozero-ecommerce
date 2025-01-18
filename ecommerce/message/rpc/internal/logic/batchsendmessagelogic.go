package logic

import (
	"context"
	"database/sql"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
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

	// Send messages through RMQ
	for _, msg := range messages {
		resp.MessageIds = append(resp.MessageIds, int64(msg.Id))

		// Create message created event
		event := &types.MessageCreatedData{
			ID:          int64(msg.Id),
			UserID:      int64(msg.UserId),
			Title:       msg.Title,
			Content:     msg.Content,
			Type:        int32(msg.Type),
			SendChannel: int32(msg.SendChannel),
			ExtraData:   msg.ExtraData.String,
		}

		// Publish to RMQ
		metadata := types.Metadata{
			Source:  "message-service",
			UserID:  int64(msg.UserId),
			TraceID: uuid.New().String(),
		}

		err = l.svcCtx.Producer.PublishMessageCreated(l.ctx, event, metadata)
		if err != nil {
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
