package logic

import (
	"context"
	"database/sql"
	"flag"
	"testing"
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

type UserCouponsTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *UserCouponsLogic
}

func TestUserCouponsSuite(t *testing.T) {
    suite.Run(t, new(UserCouponsTestSuite))
}

func (s *UserCouponsTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/marketing.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewUserCouponsLogic(context.Background(), s.ctx)
}

func (s *UserCouponsTestSuite) TestUserCoupons() {
    tests := []struct {
        name    string
        setup   func() error
        req     *marketing.UserCouponsRequest
        wantErr error
        check   func(*marketing.UserCouponsResponse)
    }{
        {
            name: "list user coupons",
            setup: func() error {
                now := time.Now()
                // Create test coupon
                couponResult, err := s.ctx.CouponsModel.Insert(context.Background(), &model.Coupons{
                    Name:      "Test Coupon",
                    Code:      "TEST" + time.Now().Format("150405"),
                    Type:      1,
                    Value:     20,
                    MinAmount: 100,
                    Status:    1,
                    StartTime: sql.NullTime{Time: now.Add(-time.Hour), Valid: true},
                    EndTime:   sql.NullTime{Time: now.Add(time.Hour), Valid: true},
                    Total:     100,
                    Received:  1,
                    Used:      0,
                })
                if err != nil {
                    return err
                }

                couponId, err := couponResult.LastInsertId()
                if err != nil {
                    return err
                }

                // Create user coupon
                _, err = s.ctx.UserCouponsModel.Insert(context.Background(), &model.UserCoupons{
                    UserId:   1001,
                    CouponId: uint64(couponId),
                    Status:   0,
                })
                return err
            },
            req: &marketing.UserCouponsRequest{
                UserId:   1001,
                Status:   0,
                Page:     1,
                PageSize: 10,
            },
            check: func(resp *marketing.UserCouponsResponse) {
                assert.Equal(s.T(), int64(1), resp.Total)
                assert.Len(s.T(), resp.Coupons, 1)
                assert.Equal(s.T(), "Test Coupon", resp.Coupons[0].Coupon.Name)
                assert.Equal(s.T(), int32(0), resp.Coupons[0].Status)
            },
        },
        {
            name: "empty result",
            req: &marketing.UserCouponsRequest{
                UserId:   9999,
                Status:   0,
                Page:     1,
                PageSize: 10,
            },
            check: func(resp *marketing.UserCouponsResponse) {
                assert.Equal(s.T(), int64(0), resp.Total)
                assert.Empty(s.T(), resp.Coupons)
            },
        },
        {
            name: "invalid user id",
            req: &marketing.UserCouponsRequest{
                UserId: 0,
            },
            wantErr: zeroerr.ErrInvalidMarketingParam,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            if tt.setup != nil {
                err := tt.setup()
                assert.NoError(t, err)
            }

            resp, err := s.logic.UserCoupons(tt.req)
            if tt.wantErr != nil {
                assert.Equal(t, tt.wantErr, err)
                return
            }

            assert.NoError(t, err)
            if tt.check != nil {
                tt.check(resp)
            }
        })
    }
}