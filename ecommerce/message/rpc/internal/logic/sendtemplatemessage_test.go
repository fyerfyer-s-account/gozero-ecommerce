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

type SendTemplateMessageTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *SendTemplateMessageLogic
}

func TestSendTemplateMessageSuite(t *testing.T) {
    suite.Run(t, new(SendTemplateMessageTestSuite))
}

func (s *SendTemplateMessageTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/message.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewSendTemplateMessageLogic(context.Background(), s.ctx)
}

func (s *SendTemplateMessageTestSuite) SetupTest() {
    s.cleanData()
}

func (s *SendTemplateMessageTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *SendTemplateMessageTestSuite) cleanData() {
    templates, _ := s.ctx.MessageTemplatesModel.FindByCode(context.Background(), "TEST_TPL")
    if templates != nil {
        _ = s.ctx.MessageTemplatesModel.Delete(context.Background(), templates.Id)
    }
}

func (s *SendTemplateMessageTestSuite) TestSendTemplateMessage() {
    tests := []struct {
        name    string
        setup   func() error
        req     *message.SendTemplateMessageRequest
        wantErr error
        check   func(resp *message.SendTemplateMessageResponse)
    }{
        {
            name: "send to all channels",
            setup: func() error {
                template := &model.MessageTemplates{
                    Code:            "TEST_TPL",
                    Name:            "Test Template",
                    TitleTemplate:   "Hello ${name}",
                    ContentTemplate: "Welcome ${name} to our platform!",
                    Type:            1,
                    Channels:        "[1,2,3]",
                    Status:          1,
                }
                _, err := s.ctx.MessageTemplatesModel.Insert(context.Background(), template)
                return err
            },
            req: &message.SendTemplateMessageRequest{
                TemplateCode: "TEST_TPL",
                UserId:      1,
                Params: map[string]string{
                    "name": "test user",
                },
                Channels: []int32{1, 2, 3},
            },
            wantErr: nil,
            check: func(resp *message.SendTemplateMessageResponse) {
                assert.Len(s.T(), resp.MessageIds, 3)
                // Verify messages
                for i, msgId := range resp.MessageIds {
                    msg, err := s.ctx.MessagesModel.FindOne(context.Background(), uint64(msgId))
                    assert.NoError(s.T(), err)
                    assert.Equal(s.T(), "Hello test user", msg.Title)
                    assert.Equal(s.T(), "Welcome test user to our platform!", msg.Content)
                    assert.Equal(s.T(), int64(i+1), msg.SendChannel)
                }
            },
        },
        {
            name: "template not found",
            setup: func() error { return nil },
            req: &message.SendTemplateMessageRequest{
                TemplateCode: "NON_EXISTENT",
                UserId:      1,
                Params:      map[string]string{},
                Channels:    []int32{1},
            },
            wantErr: zeroerr.ErrTemplateNotFound,
        },
        {
            name: "invalid channel",
            setup: func() error {
                template := &model.MessageTemplates{
                    Code:            "TEST_TPL",
                    Name:            "Test Template",
                    TitleTemplate:   "Hello ${name}",
                    ContentTemplate: "Welcome ${name}!",
                    Type:            1,
                    Channels:        "[1]",
                    Status:          1,
                }
                _, err := s.ctx.MessageTemplatesModel.Insert(context.Background(), template)
                return err
            },
            req: &message.SendTemplateMessageRequest{
                TemplateCode: "TEST_TPL",
                UserId:      1,
                Params:      map[string]string{"name": "test"},
                Channels:    []int32{2}, // Not allowed channel
            },
            wantErr: zeroerr.ErrInvalidSendChannel,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            s.cleanData()
            err := tt.setup()
            assert.NoError(t, err)

            resp, err := s.logic.SendTemplateMessage(tt.req)
            if tt.wantErr != nil {
                assert.Equal(t, tt.wantErr, err)
                return
            }

            assert.NoError(t, err)
            assert.NotNil(t, resp)
            tt.check(resp)
        })
    }
}