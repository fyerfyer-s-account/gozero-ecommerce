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

type ListMessagesTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *ListMessagesLogic
}

func TestListMessagesSuite(t *testing.T) {
    suite.Run(t, new(ListMessagesTestSuite))
}

func (s *ListMessagesTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/message.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewListMessagesLogic(context.Background(), s.ctx)
}

func (s *ListMessagesTestSuite) SetupTest() {
    s.cleanData()
}

func (s *ListMessagesTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *ListMessagesTestSuite) cleanData() {
    // Clean test data for user 1
    messages, _ := s.ctx.MessagesModel.FindByUserId(context.Background(), 1, 0, false, 1, 100)
    for _, msg := range messages {
        _ = s.ctx.MessagesModel.Delete(context.Background(), msg.Id)
    }
}

func (s *ListMessagesTestSuite) TestListMessages() {
    tests := []struct {
        name    string
        setup   func()
        req     *message.ListMessagesRequest
        wantErr bool
        check   func(*message.ListMessagesResponse)
    }{
        {
            name: "normal case",
            setup: func() {
                msgs := []*model.Messages{
                    {
                        UserId:  1,
                        Title:   "Test Message 1",
                        Content: "Test Content 1",
                        Type:    1,
                        ExtraData: sql.NullString{
                            String: `{"key":"value1"}`,
                            Valid:  true,
                        },
                        IsRead: 1,
                        ReadTime: sql.NullTime{
                            Time:  time.Now(),
                            Valid: true,
                        },
                    },
                    {
                        UserId:  1,
                        Title:   "Test Message 2",
                        Content: "Test Content 2",
                        Type:    1,
                        ExtraData: sql.NullString{
                            String: `{"key":"value2"}`,
                            Valid:  true,
                        },
                        IsRead: 0,
                    },
                }
                for _, msg := range msgs {
                    s.ctx.MessagesModel.Insert(context.Background(), msg)
                }
            },
            req: &message.ListMessagesRequest{
                UserId:   1,
                Page:     1,
                PageSize: 10,
            },
            wantErr: false,
            check: func(resp *message.ListMessagesResponse) {
                assert.Equal(s.T(), int64(2), resp.Total)
                assert.Len(s.T(), resp.Messages, 2)
                assert.Equal(s.T(), "Test Message 1", resp.Messages[0].Title)
                assert.Equal(s.T(), "Test Message 2", resp.Messages[1].Title)
                assert.True(s.T(), resp.Messages[0].IsRead)
                assert.False(s.T(), resp.Messages[1].IsRead)
            },
        },
        {
            name: "filter by type",
            setup: func() {
                msgs := []*model.Messages{
                    {
                        UserId:  1,
                        Title:   "System Message",
                        Content: "System Content",
                        Type:    1,
                    },
                    {
                        UserId:  1,
                        Title:   "Order Message",
                        Content: "Order Content",
                        Type:    2,
                    },
                }
                for _, msg := range msgs {
                    s.ctx.MessagesModel.Insert(context.Background(), msg)
                }
            },
            req: &message.ListMessagesRequest{
                UserId:   1,
                Type:     1,
                Page:     1,
                PageSize: 10,
            },
            wantErr: false,
            check: func(resp *message.ListMessagesResponse) {
                assert.Equal(s.T(), int64(1), resp.Total)
                assert.Len(s.T(), resp.Messages, 1)
                assert.Equal(s.T(), "System Message", resp.Messages[0].Title)
            },
        },
        {
            name: "unread only",
            setup: func() {
                msgs := []*model.Messages{
                    {
                        UserId:  1,
                        Title:   "Read Message",
                        Content: "Read Content",
                        IsRead:  1,
                    },
                    {
                        UserId:  1,
                        Title:   "Unread Message",
                        Content: "Unread Content",
                        IsRead:  0,
                    },
                }
                for _, msg := range msgs {
                    s.ctx.MessagesModel.Insert(context.Background(), msg)
                }
            },
            req: &message.ListMessagesRequest{
                UserId:     1,
                UnreadOnly: true,
                Page:       1,
                PageSize:   10,
            },
            wantErr: false,
            check: func(resp *message.ListMessagesResponse) {
                assert.Equal(s.T(), int64(1), resp.Total)
                assert.Len(s.T(), resp.Messages, 1)
                assert.Equal(s.T(), "Unread Message", resp.Messages[0].Title)
            },
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            s.cleanData()
            tt.setup()

            resp, err := s.logic.ListMessages(tt.req)
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