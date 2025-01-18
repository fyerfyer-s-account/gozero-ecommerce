package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SendTemplateMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendTemplateMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendTemplateMessageLogic {
	return &SendTemplateMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendTemplateMessageLogic) SendTemplateMessage(in *message.SendTemplateMessageRequest) (*message.SendTemplateMessageResponse, error) {
	// Find template
	template, err := l.svcCtx.MessageTemplatesModel.FindByCode(l.ctx, in.TemplateCode)
	if err != nil {
		return nil, zeroerr.ErrTemplateNotFound
	}

	// Parse channels
	var templateChannels []int32
	err = json.Unmarshal([]byte(template.Channels), &templateChannels)
	if err != nil {
		return nil, zeroerr.ErrInvalidTemplate
	}

	// Validate channels
	selectedChannels := make([]int32, 0)
	for _, ch := range in.Channels {
		if contains(templateChannels, ch) {
			selectedChannels = append(selectedChannels, ch)
		}
	}
	if len(selectedChannels) == 0 {
		return nil, zeroerr.ErrInvalidSendChannel
	}

	// Parse template with params
	title := parseTemplate(template.TitleTemplate, in.Params)
	content := parseTemplate(template.ContentTemplate, in.Params)

	var messageIds []int64
	err = l.svcCtx.MessagesModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// Create message for each channel
		for _, channel := range selectedChannels {
			// Create message
			msg := &model.Messages{
				UserId:      uint64(in.UserId),
				Title:       title,
				Content:     content,
				Type:        template.Type,
				SendChannel: int64(channel),
				ExtraData: sql.NullString{
					String: getExtraData(in.Params),
					Valid:  true,
				},
				IsRead: 0,
			}

			result, err := l.svcCtx.MessagesModel.Insert(ctx, msg)
			if err != nil {
				return zeroerr.ErrMessageCreateFailed
			}

			messageId, _ := result.LastInsertId()
			messageIds = append(messageIds, messageId)

			// Create message send record
			send := &model.MessageSends{
				MessageId: uint64(messageId),
				TemplateId: sql.NullInt64{
					Int64: int64(template.Id),
					Valid: true,
				},
				UserId:     uint64(in.UserId),
				Channel:    int64(channel),
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
				Title:       title,
				Content:     content,
				Type:        int32(template.Type),
				SendChannel: channel,
				ExtraData:   getExtraData(in.Params),
			}, types.Metadata{
				Source: "message-service",
				UserID: in.UserId,
			})
			if err != nil {
				logx.Errorf("Failed to publish message created event: %v", err)
				// Continue since message is created
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &message.SendTemplateMessageResponse{
		MessageIds: messageIds,
	}, nil
}

func contains(arr []int32, val int32) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func parseTemplate(template string, params map[string]string) string {
	result := template
	for key, value := range params {
		result = strings.ReplaceAll(result, "${"+key+"}", value)
	}
	return result
}

func getExtraData(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	data, _ := json.Marshal(params)
	return string(data)
}
