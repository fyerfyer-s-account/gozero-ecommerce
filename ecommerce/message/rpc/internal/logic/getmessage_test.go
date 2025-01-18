package logic

import (
    "context"
    "database/sql"
    "flag"
    "testing"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/zeromicro/go-zero/core/conf"
)

type GetMessageTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *GetMessageLogic
}

func TestGetMessageSuite(t *testing.T) {
    suite.Run(t, new(GetMessageTestSuite))
}

func (s *GetMessageTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/message.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewGetMessageLogic(context.Background(), s.ctx)
}

func (s *GetMessageTestSuite) SetupTest() {
    s.cleanData()
}

func (s *GetMessageTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *GetMessageTestSuite) cleanData() {
    // Clean test data
    messages, _ := s.ctx.MessagesModel.FindOne(context.Background(), 1)
    if messages != nil {
        _ = s.ctx.MessagesModel.Delete(context.Background(), messages.Id)
    }
}

func (s *GetMessageTestSuite) TestGetMessage() {
    tests := []struct {
        name    string
        setup   func() int64
        req     *message.GetMessageRequest
        wantErr bool
        check   func(*message.GetMessageResponse)
    }{
        {
            name: "normal case",
            setup: func() int64 {
                msg := &model.Messages{
                    UserId:      1,
                    Title:       "Test Message",
                    Content:     "Test Content",
                    Type:        1,
                    SendChannel: 1,
                    ExtraData: sql.NullString{
                        String: `{"key":"value"}`,
                        Valid:  true,
                    },
                    IsRead: 1,
                    ReadTime: sql.NullTime{
                        Time:  time.Now(),
                        Valid: true,
                    },
                }
                result, _ := s.ctx.MessagesModel.Insert(context.Background(), msg)
                id, _ := result.LastInsertId()
                return id
            },
            req: &message.GetMessageRequest{
                MessageId: 1,
            },
            wantErr: false,
            check: func(resp *message.GetMessageResponse) {
				assert.NotNil(s.T(), resp.Message)
				assert.Equal(s.T(), "Test Message", resp.Message.Title)
				assert.Equal(s.T(), "Test Content", resp.Message.Content)
				assert.Equal(s.T(), int32(1), resp.Message.Type)
				assert.Equal(s.T(), int32(1), resp.Message.SendChannel)
				assert.JSONEq(s.T(), `{"key":"value"}`, resp.Message.ExtraData) // 使用 JSONEq
				assert.True(s.T(), resp.Message.IsRead)
			},
        },
        {
            name:  "message not found",
            setup: func() int64 { return 0 },
            req: &message.GetMessageRequest{
                MessageId: 999,
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

            resp, err := s.logic.GetMessage(tt.req)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)
            assert.NotNil(t, resp)
            tt.check(resp)
        })
    }
}