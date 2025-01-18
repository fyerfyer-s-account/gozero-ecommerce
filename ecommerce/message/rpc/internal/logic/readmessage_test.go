package logic

import (
	"context"
	"flag"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/core/conf"
)

type ReadMessageTestSuite struct {
	suite.Suite
	ctx   *svc.ServiceContext
	logic *ReadMessageLogic
}

func TestReadMessageSuite(t *testing.T) {
	suite.Run(t, new(ReadMessageTestSuite))
}

func (s *ReadMessageTestSuite) SetupSuite() {
	configFile := flag.String("f", "../../etc/message.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
	s.logic = NewReadMessageLogic(context.Background(), s.ctx)
}

func (s *ReadMessageTestSuite) SetupTest() {
	s.cleanData()
}

func (s *ReadMessageTestSuite) TearDownTest() {
	s.cleanData()
}

func (s *ReadMessageTestSuite) cleanData() {
	messages, _ := s.ctx.MessagesModel.FindOne(context.Background(), 1)
	if messages != nil {
		_ = s.ctx.MessagesModel.Delete(context.Background(), messages.Id)
	}
}

func (s *ReadMessageTestSuite) TestReadMessage() {
	tests := []struct {
		name    string
		setup   func() uint64
		req     *message.ReadMessageRequest
		wantErr error
	}{
		{
			name: "read unread message",
			setup: func() uint64 {
				msg := &model.Messages{
					UserId:  1,
					Title:   "Test Message",
					Content: "Test Content",
					Type:    1,
					IsRead:  0,
				}
				result, _ := s.ctx.MessagesModel.Insert(context.Background(), msg)
				id, _ := result.LastInsertId()
				return uint64(id)
			},
			req: &message.ReadMessageRequest{
				MessageId: 1,
				UserId:    1,
			},
			wantErr: nil,
		},
		{
			name:  "message not found",
			setup: func() uint64 { return 0 },
			req: &message.ReadMessageRequest{
				MessageId: 999,
				UserId:    1,
			},
			wantErr: zeroerr.ErrMessageNotFound,
		},
		{
			name: "unauthorized user",
			setup: func() uint64 {
				msg := &model.Messages{
					UserId:  1,
					Title:   "Test Message",
					Content: "Test Content",
					Type:    1,
					IsRead:  0,
				}
				result, _ := s.ctx.MessagesModel.Insert(context.Background(), msg)
				id, _ := result.LastInsertId()
				return uint64(id)
			},
			req: &message.ReadMessageRequest{
				MessageId: 1,
				UserId:    2,
			},
			wantErr: zeroerr.ErrMessageNotFound,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s.cleanData()
			msgId := tt.setup()
			if msgId > 0 {
				tt.req.MessageId = int64(msgId)
			}

			resp, err := s.logic.ReadMessage(tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}

			assert.NoError(t, err)
			assert.True(t, resp.Success)

			// Verify message is marked as read
			msg, err := s.ctx.MessagesModel.FindOne(context.Background(), msgId)
			assert.NoError(t, err)
			assert.Equal(t, int64(1), msg.IsRead)
			assert.True(t, msg.ReadTime.Valid)
		})
	}
}
