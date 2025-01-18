package logic

import (
	"context"
	"flag"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/core/conf"
)

type BatchSendMessageTestSuite struct {
	suite.Suite
	ctx   *svc.ServiceContext
	logic *BatchSendMessageLogic
}

func TestBatchSendMessageSuite(t *testing.T) {
	suite.Run(t, new(BatchSendMessageTestSuite))
}

func (s *BatchSendMessageTestSuite) SetupSuite() {
	configFile := flag.String("f", "../../etc/message.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
	s.logic = NewBatchSendMessageLogic(context.Background(), s.ctx)
}

func (s *BatchSendMessageTestSuite) SetupTest() {
	s.cleanData()
}

func (s *BatchSendMessageTestSuite) TearDownTest() {
	s.cleanData()
}

func (s *BatchSendMessageTestSuite) cleanData() {
	ctx := context.Background()
	// Clean test messages - get all types, including read messages
	messages, _ := s.ctx.MessagesModel.FindByUserId(ctx, 1001, 0, false, 1, 100)
	for _, msg := range messages {
		_ = s.ctx.MessagesModel.Delete(ctx, msg.Id)
	}
}

func (s *BatchSendMessageTestSuite) TestBatchSendMessage() {
	tests := []struct {
		name  string
		req   *message.BatchSendMessageRequest
		check func(*message.BatchSendMessageResponse)
	}{
		{
			name: "send to multiple users",
			req: &message.BatchSendMessageRequest{
				UserIds:     []int64{1001, 1002},
				Title:       "Test Message",
				Content:     "Test Content",
				Type:        types.MessageTypeSystem,
				SendChannel: types.ChannelInApp,
				ExtraData:   `{"key":"value"}`,
			},
			check: func(resp *message.BatchSendMessageResponse) {
				assert.Len(s.T(), resp.MessageIds, 2)
				assert.Empty(s.T(), resp.Errors)

				// Verify messages in DB - get system messages only
				ctx := context.Background()
				messages, err := s.ctx.MessagesModel.FindByUserId(ctx, 1001, types.MessageTypeSystem, false, 1, 10)
				assert.NoError(s.T(), err)
				assert.Len(s.T(), messages, 1)
				assert.Equal(s.T(), "Test Message", messages[0].Title)
			},
		},
		{
			name: "empty user list",
			req: &message.BatchSendMessageRequest{
				UserIds: []int64{},
			},
			check: func(resp *message.BatchSendMessageResponse) {
				assert.Empty(s.T(), resp.MessageIds)
				assert.Empty(s.T(), resp.Errors)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			resp, err := s.logic.BatchSendMessage(tt.req)
			assert.NoError(t, err)
			tt.check(resp)
		})
	}
}
