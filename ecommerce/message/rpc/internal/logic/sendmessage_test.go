package logic

import (
	"context"
	"flag"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/core/conf"
)

type SendMessageTestSuite struct {
	suite.Suite
	ctx   *svc.ServiceContext
	logic *SendMessageLogic
}

func TestSendMessageSuite(t *testing.T) {
	suite.Run(t, new(SendMessageTestSuite))
}

func (s *SendMessageTestSuite) SetupSuite() {
	configFile := flag.String("f", "../../etc/message.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
	s.logic = NewSendMessageLogic(context.Background(), s.ctx)
}

func (s *SendMessageTestSuite) SetupTest() {
	s.cleanData()
}

func (s *SendMessageTestSuite) TearDownTest() {
	s.cleanData()
}

func (s *SendMessageTestSuite) cleanData() {
	messages, _ := s.ctx.MessagesModel.FindByUserId(context.Background(), 1, 0, false, 1, 100)
	for _, msg := range messages {
		_ = s.ctx.MessageSendsModel.Delete(context.Background(), msg.Id)
		_ = s.ctx.MessagesModel.Delete(context.Background(), msg.Id)
	}
}

func (s *SendMessageTestSuite) TestSendMessage() {
	tests := []struct {
		name    string
		req     *message.SendMessageRequest
		wantErr error
	}{
		{
			name: "normal case",
			req: &message.SendMessageRequest{
				UserId:      1,
				Title:       "Test Message",
				Content:     "Test Content",
				Type:        1,
				SendChannel: 1,
				ExtraData:   `{"key":"value"}`,
			},
			wantErr: nil,
		},
		{
			name: "empty title",
			req: &message.SendMessageRequest{
				UserId:      1,
				Title:       "",
				Content:     "Test Content",
				Type:        1,
				SendChannel: 1,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "invalid message type",
			req: &message.SendMessageRequest{
				UserId:      1,
				Title:       "Test Message",
				Content:     "Test Content",
				Type:        0,
				SendChannel: 1,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "invalid send channel",
			req: &message.SendMessageRequest{
				UserId:      1,
				Title:       "Test Message",
				Content:     "Test Content",
				Type:        1,
				SendChannel: 5,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s.cleanData()

			resp, err := s.logic.SendMessage(tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, resp.MessageId)

			// 验证消息是否正确创建
			msg, err := s.ctx.MessagesModel.FindOne(context.Background(), uint64(resp.MessageId))
			assert.NoError(t, err)
			assert.NotNil(t, msg)
			assert.Equal(t, tt.req.Title, msg.Title)
			assert.Equal(t, tt.req.Content, msg.Content)
			assert.Equal(t, int64(tt.req.Type), msg.Type)
			assert.Equal(t, int64(tt.req.SendChannel), msg.SendChannel)

			sendRecords, _ := s.ctx.MessageSendsModel.FindByMessageId(context.Background(), uint64(resp.MessageId))
			assert.NotEmpty(t, sendRecords)
			assert.Equal(t, uint64(resp.MessageId), sendRecords[0].MessageId)
			assert.Equal(t, uint64(tt.req.UserId), sendRecords[0].UserId)
			assert.Equal(t, int64(tt.req.SendChannel), sendRecords[0].Channel)
		})
	}
}
