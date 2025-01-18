package logic

import (
	"context"
	"database/sql"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 消息发送
func (l *SendMessageLogic) SendMessage(in *message.SendMessageRequest) (*message.SendMessageResponse, error) {
	// Validate input
	if in.Title == "" || in.Content == "" {
		return nil, zeroerr.ErrInvalidParam
	}
	if !isValidMessageType(in.Type) || !isValidSendChannel(in.SendChannel) {
		return nil, zeroerr.ErrInvalidParam
	}

	var messageId int64
	err := l.svcCtx.MessagesModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// Create message
		msg := &model.Messages{
			UserId:      uint64(in.UserId),
			Title:       in.Title,
			Content:     in.Content,
			Type:        int64(in.Type),
			SendChannel: int64(in.SendChannel),
			ExtraData: sql.NullString{
				String: in.ExtraData,
				Valid:  in.ExtraData != "",
			},
			IsRead: 0,
		}
		result, err := l.svcCtx.MessagesModel.Insert(ctx, msg)
		if err != nil {
			return zeroerr.ErrMessageCreateFailed
		}

		messageId, err = result.LastInsertId()
		if err != nil {
			return zeroerr.ErrMessageCreateFailed
		}

		// Create message send record
		send := &model.MessageSends{
			MessageId:  uint64(messageId),
			UserId:     uint64(in.UserId),
			Channel:    int64(in.SendChannel),
			Status:     1, // Pending
			RetryCount: 0,
		}
		_, err = l.svcCtx.MessageSendsModel.Insert(ctx, send)
		if err != nil {
			return zeroerr.ErrMessageCreateFailed
		}

		// Publish message created event
		err = l.svcCtx.Producer.PublishMessageCreated(ctx, &types.MessageCreatedData{
			ID:          messageId,
			UserID:      in.UserId,
			Title:       in.Title,
			Content:     in.Content,
			Type:        in.Type,
			SendChannel: in.SendChannel,
			ExtraData:   in.ExtraData,
		}, types.Metadata{
			Source: "message-service",
			UserID: in.UserId,
		})
		if err != nil {
			logx.Errorf("Failed to publish message created event: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &message.SendMessageResponse{
		MessageId: messageId, // Return actual messageId
	}, nil
}

func isValidMessageType(t int32) bool {
	return t >= 1 && t <= 4
}

func isValidSendChannel(c int32) bool {
	return c >= 1 && c <= 4
}
