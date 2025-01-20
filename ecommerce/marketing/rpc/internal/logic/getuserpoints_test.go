package logic

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GetUserPointsTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *GetUserPointsLogic
}

func TestGetUserPointsSuite(t *testing.T) {
    suite.Run(t, new(GetUserPointsTestSuite))
}

func (s *GetUserPointsTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/marketing.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
    s.logic = NewGetUserPointsLogic(context.Background(), s.ctx)
}

func (s *GetUserPointsTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *GetUserPointsTestSuite) cleanData() {
    _ = s.ctx.UserPointsModel.Delete(context.Background(), 1001)
    _ = s.ctx.UserPointsModel.Delete(context.Background(), 1002)
}

func (s *GetUserPointsTestSuite) TestGetUserPoints() {
    tests := []struct {
        name    string
        setup   func() error
        req     *marketing.GetUserPointsRequest
        want    *marketing.GetUserPointsResponse
        wantErr error
    }{
        {
            name: "get existing user points",
            setup: func() error {
                return s.ctx.UserPointsModel.Trans(context.Background(), func(ctx context.Context, session sqlx.Session) error {
                    _, err := s.ctx.UserPointsModel.Insert(ctx, &model.UserPoints{
                        UserId:      1001,
                        Points:      100,
                        TotalPoints: 150,
                        UsedPoints:  50,
                        UpdatedAt:   time.Now(),
                    })
                    return err
                })
            },
            req: &marketing.GetUserPointsRequest{
                UserId: 1001,
            },
            want: &marketing.GetUserPointsResponse{
                Points: 100,
            },
            wantErr: nil,
        },
        {
            name: "get non-existent user points",
            req: &marketing.GetUserPointsRequest{
                UserId: 9999,
            },
            want: &marketing.GetUserPointsResponse{
                Points: 0,
            },
            wantErr: nil,
        },
        {
            name: "invalid user id",
            req: &marketing.GetUserPointsRequest{
                UserId: 0,
            },
            wantErr: zeroerr.ErrInvalidMarketingParam,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            s.cleanData()
            if tt.setup != nil {
                err := tt.setup()
                assert.NoError(t, err)
            }

            resp, err := s.logic.GetUserPoints(tt.req)
            if tt.wantErr != nil {
                assert.Equal(t, tt.wantErr, err)
                return
            }

            assert.NoError(t, err)
            assert.Equal(t, tt.want.Points, resp.Points)
        })
    }
}