package logic

import (
    "context"
    "flag"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
    "github.com/zeromicro/go-zero/core/conf"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type UsePointsTestSuite struct {
    suite.Suite
    ctx    context.Context
    cancel context.CancelFunc
    svcCtx *svc.ServiceContext
    logic  *UsePointsLogic
}

func (s *UsePointsTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/marketing.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    // Increase timeout to 30s
    s.ctx, s.cancel = context.WithTimeout(context.Background(), 30*time.Second)
    s.svcCtx = svc.NewServiceContext(c)
    s.logic = NewUsePointsLogic(s.ctx, s.svcCtx)
}

func (s *UsePointsTestSuite) cleanData() {
    // Clean user points
    _ = s.svcCtx.UserPointsModel.Delete(s.ctx, 1001)
    _ = s.svcCtx.UserPointsModel.Delete(s.ctx, 1002)
    
    // Clean points records
    records, err := s.svcCtx.PointsRecordsModel.FindByUserId(s.ctx, 1001, 1, 100)
    if err == nil {
        for _, r := range records {
            _ = s.svcCtx.PointsRecordsModel.Delete(s.ctx, r.Id)
        }
    }
}

func (s *UsePointsTestSuite) TestUsePoints() {
    tests := []struct {
        name    string
        setup   func() error
        req     *marketing.UsePointsRequest
        wantErr error
        check   func(*marketing.UsePointsResponse)
    }{
        {
            name: "invalid parameters (zero user ID)",
            req: &marketing.UsePointsRequest{
                UserId:  0,
                Points:  50,
                Usage:   "order",
                OrderNo: "ORD123",
            },
            wantErr: zeroerr.ErrInvalidMarketingParam,
        },
        {
            name: "insufficient points",
            setup: func() error {
                _, err := s.svcCtx.UserPointsModel.Insert(s.ctx, &model.UserPoints{
                    UserId:      1001,
                    Points:      30,
                    TotalPoints: 30,
                    UsedPoints:  0,
                    UpdatedAt:   time.Now(),
                })
                return err
            },
            req: &marketing.UsePointsRequest{
                UserId:  1001,
                Points:  50,
                Usage:   "order",
                OrderNo: "ORD123",
            },
            wantErr: zeroerr.ErrInsufficientPoints,
        },
        {
            name: "successful use of points",
            setup: func() error {
                _, err := s.svcCtx.UserPointsModel.Insert(s.ctx, &model.UserPoints{
                    UserId:      1001,
                    Points:      100,
                    TotalPoints: 100,
                    UsedPoints:  0,
                    UpdatedAt:   time.Now(),
                })
                return err
            },
            req: &marketing.UsePointsRequest{
                UserId:  1001,
                Points:  50,
                Usage:   "order",
                OrderNo: "ORD123",
            },
            check: func(resp *marketing.UsePointsResponse) {
                if assert.NotNil(s.T(), resp) {
                    assert.True(s.T(), resp.Success)
                    assert.Equal(s.T(), int64(50), resp.CurrentPoints)

                    // Verify points record exists
                    records, err := s.svcCtx.PointsRecordsModel.FindByUserId(s.ctx, 1001, 1, 10)
                    assert.NoError(s.T(), err)
                    if assert.NotEmpty(s.T(), records) {
                        assert.Equal(s.T(), int64(50), records[0].Points)
                    }
                }
            },
        },
    }

    for _, tt := range tests {
        s.Run(tt.name, func() {
            s.cleanData()
            if tt.setup != nil {
                err := tt.setup()
                assert.NoError(s.T(), err)
            }

            resp, err := s.logic.UsePoints(tt.req)
            if tt.wantErr != nil {
                assert.Equal(s.T(), tt.wantErr, err)
                return
            }

            assert.NoError(s.T(), err)
            if tt.check != nil {
                tt.check(resp)
            }
        })
    }
}