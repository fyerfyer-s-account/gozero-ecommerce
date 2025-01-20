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

type ListCouponsTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *ListCouponsLogic
}

func TestListCouponsSuite(t *testing.T) {
    suite.Run(t, new(ListCouponsTestSuite))
}

func (s *ListCouponsTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/marketing.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewListCouponsLogic(context.Background(), s.ctx)
}

func (s *ListCouponsTestSuite) cleanData() {
    var ids []uint64
	
    coupons, err := s.ctx.CouponsModel.FindMany(context.Background(), 1, 1, 100)
    if err == nil {
        for _, c := range coupons {
            if len(c.Code) >= 4 && c.Code[:4] == "TEST" {
                ids = append(ids, c.Id)
            }
        }
    }

    for _, id := range ids {
        _ = s.ctx.CouponsModel.Delete(context.Background(), id)
    }
}

func (s *ListCouponsTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *ListCouponsTestSuite) TestListCoupons() {
    tests := []struct {
        name  string
        setup func() error
        req   *marketing.ListCouponsRequest
        check func(*marketing.ListCouponsResponse)
    }{
        {
            name: "list active coupons",
            setup: func() error {
                now := time.Now()
                for i := 1; i <= 3; i++ {
                    _, err := s.ctx.CouponsModel.Insert(context.Background(), &model.Coupons{
                        Name:      fmt.Sprintf("Test Coupon %d", i),
                        Code:      fmt.Sprintf("TEST_A_%d", i),
                        Type:      1,
                        Value:     float64(10 * i),
                        MinAmount: float64(50 * i),
                        Status:    1,
                        StartTime: sql.NullTime{Time: now, Valid: true},
                        EndTime:   sql.NullTime{Time: now.Add(24 * time.Hour), Valid: true},
                        Total:     100,
                        Received:  0,
                        Used:      0,
                        PerLimit:  1,
                        UserLimit: 1,
                    })
                    if err != nil {
                        return err
                    }
                }
                return nil
            },
            req: &marketing.ListCouponsRequest{
                Status:   1,
                Page:     1,
                PageSize: 10,
            },
            check: func(resp *marketing.ListCouponsResponse) {
                assert.Equal(s.T(), int64(3), resp.Total)
                assert.Len(s.T(), resp.Coupons, 3)
                for i, coupon := range resp.Coupons {
                    assert.Equal(s.T(), fmt.Sprintf("Test Coupon %d", i+1), coupon.Name)
                    assert.Equal(s.T(), int32(1), coupon.Status)
                }
            },
        },
        {
            name: "pagination test",
            setup: func() error {
                now := time.Now()
                for i := 1; i <= 5; i++ {
                    _, err := s.ctx.CouponsModel.Insert(context.Background(), &model.Coupons{
                        Name:      fmt.Sprintf("Test Coupon %d", i),
                        Code:      fmt.Sprintf("TEST_B_%d", i),
                        Type:      1,
                        Value:     float64(10 * i),
                        Status:    1,
                        StartTime: sql.NullTime{Time: now, Valid: true},
                        EndTime:   sql.NullTime{Time: now.Add(24 * time.Hour), Valid: true},
                        Total:     100,
                        Received:  0,
                        Used:      0,
                        PerLimit:  1,
                        UserLimit: 1,
                    })
                    if err != nil {
                        return err
                    }
                }
                return nil
            },
            req: &marketing.ListCouponsRequest{
                Status:   1,
                Page:     2,
                PageSize: 2,
            },
            check: func(resp *marketing.ListCouponsResponse) {
                assert.Equal(s.T(), int64(5), resp.Total)
                assert.Len(s.T(), resp.Coupons, 2)
            },
        },
        {
            name: "empty result",
            req: &marketing.ListCouponsRequest{
                Status:   2,
                Page:     1,
                PageSize: 10,
            },
            check: func(resp *marketing.ListCouponsResponse) {
                assert.Equal(s.T(), int64(0), resp.Total)
                assert.Empty(s.T(), resp.Coupons)
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

            resp, err := s.logic.ListCoupons(tt.req)
            assert.NoError(t, err)
            if tt.check != nil {
                tt.check(resp)
            }
        })
    }
}