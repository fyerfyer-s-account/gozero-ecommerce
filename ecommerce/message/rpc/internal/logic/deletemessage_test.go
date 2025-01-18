package logic

import (
	"context"
	"flag"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/core/conf"
)

type DeleteMessageTestSuite struct {
	suite.Suite
	ctx   *svc.ServiceContext
	logic *DeleteMessageLogic
}

func TestDeleteMessageSuite(t *testing.T) {
	suite.Run(t, new(DeleteMessageTestSuite))
}

func (s *DeleteMessageTestSuite) SetupSuite() {
	configFile := flag.String("f", "../../etc/message.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
	s.logic = NewDeleteMessageLogic(context.Background(), s.ctx)
}

func (s *DeleteMessageTestSuite) SetupTest() {
	s.cleanData()
}

func (s *DeleteMessageTestSuite) TearDownTest() {
	s.cleanData()
}

func (s *DeleteMessageTestSuite) cleanData() {
	// Clean test data
	messages, _ := s.ctx.MessagesModel.FindOne(context.Background(), 1)
	if messages != nil {
		_ = s.ctx.MessagesModel.Delete(context.Background(), messages.Id)
	}
}

func (s *DeleteMessageTestSuite) TestDeleteMessage() {
	tests := []struct {
		name    string
		setup   func() int64
		req     *message.DeleteMessageRequest
		wantErr bool
	}{
		{
			name: "normal case",
			setup: func() int64 {
				msg := &model.Messages{
					UserId:  1,
					Title:   "Test Message",
					Content: "Test Content",
					Type:    1,
				}
				result, _ := s.ctx.MessagesModel.Insert(context.Background(), msg)
				id, _ := result.LastInsertId()
				return id
			},
			req: &message.DeleteMessageRequest{
				MessageId: 1,
				UserId:    1,
			},
			wantErr: false,
		},
		{
			name:  "message not found",
			setup: func() int64 { return 0 },
			req: &message.DeleteMessageRequest{
				MessageId: 999,
				UserId:    1,
			},
			wantErr: true,
		},
		{
			name: "unauthorized user",
			setup: func() int64 {
				msg := &model.Messages{
					UserId:  1,
					Title:   "Test Message",
					Content: "Test Content",
					Type:    1,
				}
				result, _ := s.ctx.MessagesModel.Insert(context.Background(), msg)
				id, _ := result.LastInsertId()
				return id
			},
			req: &message.DeleteMessageRequest{
				MessageId: 1,
				UserId:    2,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			msgId := tt.setup()
			if msgId > 0 {
				tt.req.MessageId = msgId
			}

			resp, err := s.logic.DeleteMessage(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			
			assert.NoError(t, err)
			assert.NotNil(t, resp)  
			assert.True(t, resp.Success)

			// Verify message is deleted
			_, err = s.ctx.MessagesModel.FindOne(context.Background(), uint64(tt.req.MessageId))
			assert.Error(t, err)
		})
	}
}
