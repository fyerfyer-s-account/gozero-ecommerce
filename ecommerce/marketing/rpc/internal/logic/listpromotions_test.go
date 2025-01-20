package logic

import (
    "context"
    "database/sql"
    "flag"
    "fmt"
    "testing"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
    "github.com/zeromicro/go-zero/core/conf"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type ListPromotionsTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *ListPromotionsLogic
}

func TestListPromotionsSuite(t *testing.T) {
    suite.Run(t, new(ListPromotionsTestSuite))
}

func (s *ListPromotionsTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/marketing.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewListPromotionsLogic(context.Background(), s.ctx)
}

func (s *ListPromotionsTestSuite) cleanData() {
    var ids []uint64
    promos, err := s.ctx.PromotionsModel.FindByStatus(context.Background(), 1, 1, 100)
    if err == nil {
        for _, p := range promos {
            if len(p.Name) >= 4 && p.Name[:4] == "Test" {
                ids = append(ids, p.Id)
            }
        }
    }

    for _, id := range ids {
        _ = s.ctx.PromotionsModel.Delete(context.Background(), id)
    }
}

func (s *ListPromotionsTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *ListPromotionsTestSuite) TestListPromotions() {
    tests := []struct {
        name  string
        setup func() error
        req   *marketing.ListPromotionsRequest
        check func(*marketing.ListPromotionsResponse)
    }{
        {
            name: "list active promotions",
            setup: func() error {
                now := time.Now()
                for i := 1; i <= 3; i++ {
                    _, err := s.ctx.PromotionsModel.Insert(context.Background(), &model.Promotions{
                        Name:      fmt.Sprintf("Test Promotion %d", i),
                        Type:      1,
                        Rules:     fmt.Sprintf(`{"minAmount":%d,"discountAmount":%d}`, 100*i, 10*i),
                        Status:    1,
                        StartTime: sql.NullTime{Time: now, Valid: true},
                        EndTime:   sql.NullTime{Time: now.Add(24 * time.Hour), Valid: true},
                    })
                    if err != nil {
                        return err
                    }
                }
                return nil
            },
            req: &marketing.ListPromotionsRequest{
                Status:   1,
                Page:     1,
                PageSize: 10,
            },
            check: func(resp *marketing.ListPromotionsResponse) {
                assert.Equal(s.T(), int64(3), resp.Total)
                assert.Len(s.T(), resp.Promotions, 3)
                for i, promo := range resp.Promotions {
                    assert.Equal(s.T(), fmt.Sprintf("Test Promotion %d", i+1), promo.Name)
                    assert.Equal(s.T(), int32(1), promo.Status)
                }
            },
        },
        {
            name: "pagination test",
            setup: func() error {
                now := time.Now()
                for i := 1; i <= 5; i++ {
                    _, err := s.ctx.PromotionsModel.Insert(context.Background(), &model.Promotions{
                        Name:      fmt.Sprintf("Test Promotion %d", i),
                        Type:      1,
                        Rules:     `{"minAmount":100,"discountAmount":10}`,
                        Status:    1,
                        StartTime: sql.NullTime{Time: now, Valid: true},
                        EndTime:   sql.NullTime{Time: now.Add(24 * time.Hour), Valid: true},
                    })
                    if err != nil {
                        return err
                    }
                }
                return nil
            },
            req: &marketing.ListPromotionsRequest{
                Status:   1,
                Page:     2,
                PageSize: 2,
            },
            check: func(resp *marketing.ListPromotionsResponse) {
                assert.Equal(s.T(), int64(5), resp.Total)
                assert.Len(s.T(), resp.Promotions, 2)
            },
        },
        {
            name: "empty result",
            req: &marketing.ListPromotionsRequest{
                Status:   2,
                Page:     1,
                PageSize: 10,
            },
            check: func(resp *marketing.ListPromotionsResponse) {
                assert.Equal(s.T(), int64(0), resp.Total)
                assert.Empty(s.T(), resp.Promotions)
            },
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            s.cleanData()
            if tt.setup != nil {
                err := tt.setup()
                assert.NoError(t, err)
            }

            resp, err := s.logic.ListPromotions(tt.req)
            assert.NoError(t, err)
            if tt.check != nil {
                tt.check(resp)
            }
        })
    }
}