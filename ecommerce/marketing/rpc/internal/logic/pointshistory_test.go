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

type PointsHistoryTestSuite struct {
	suite.Suite
	ctx   *svc.ServiceContext
	logic *PointsHistoryLogic
}

func TestPointsHistorySuite(t *testing.T) {
	suite.Run(t, new(PointsHistoryTestSuite))
}

func (s *PointsHistoryTestSuite) SetupSuite() {
	configFile := flag.String("f", "../../etc/marketing.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
	s.logic = NewPointsHistoryLogic(context.Background(), s.ctx)
}

func (s *PointsHistoryTestSuite) TearDownTest() {
	s.cleanData()
}

func (s *PointsHistoryTestSuite) cleanData() {
	var ids []uint64
	records, err := s.ctx.PointsRecordsModel.FindByUserId(context.Background(), 1001, 1, 100)
	if err == nil {
		for _, r := range records {
			ids = append(ids, r.Id)
		}
	}

	for _, id := range ids {
		_ = s.ctx.PointsRecordsModel.Delete(context.Background(), id)
	}
}

func (s *PointsHistoryTestSuite) TestPointsHistory() {
	tests := []struct {
		name    string
		setup   func() error
		req     *marketing.PointsHistoryRequest
		wantErr error
		check   func(*marketing.PointsHistoryResponse)
	}{
		{
			name: "get points history",
			setup: func() error {
				now := time.Now()
				records := []*model.PointsRecords{
					{
						UserId:    1001,
						Points:    100,
						Type:      1,
						Source:    "purchase",
						Remark:    sql.NullString{String: "Order reward", Valid: true},
						CreatedAt: now, // Most recent record
					},
					{
						UserId:    1001,
						Points:    50,
						Type:      2,
						Source:    "redeem",
						Remark:    sql.NullString{String: "Product redemption", Valid: true},
						CreatedAt: now.Add(-time.Hour),
					},
				}
				// Make sure we insert records in order
				for _, record := range records {
					if _, err := s.ctx.PointsRecordsModel.Insert(context.Background(), record); err != nil {
						return err
					}
				}
				return nil
			},
			req: &marketing.PointsHistoryRequest{
				UserId:   1001,
				Page:     1,
				PageSize: 10,
			},
			check: func(resp *marketing.PointsHistoryResponse) {
				assert.Equal(s.T(), int64(2), resp.Total)
				assert.Len(s.T(), resp.Records, 2)
				assert.Equal(s.T(), int64(100), resp.Records[0].Points)
				assert.Equal(s.T(), int32(1), resp.Records[0].Type)
				assert.Equal(s.T(), "purchase", resp.Records[0].Source)
			},
		},
		{
			name: "pagination test",
			setup: func() error {
				now := time.Now()
				var records []*model.PointsRecords
				for i := 0; i < 5; i++ {
					records = append(records, &model.PointsRecords{
						UserId:    1001,
						Points:    int64(10 * (i + 1)),
						Type:      1,
						Source:    "test",
						CreatedAt: now.Add(time.Duration(-i) * time.Hour),
					})
				}
				return s.ctx.PointsRecordsModel.BatchInsert(context.Background(), records)
			},
			req: &marketing.PointsHistoryRequest{
				UserId:   1001,
				Page:     2,
				PageSize: 2,
			},
			check: func(resp *marketing.PointsHistoryResponse) {
				assert.Equal(s.T(), int64(5), resp.Total)
				assert.Len(s.T(), resp.Records, 2)
			},
		},
		{
			name: "invalid user id",
			req: &marketing.PointsHistoryRequest{
				UserId:   0,
				Page:     1,
				PageSize: 10,
			},
			wantErr: zeroerr.ErrInvalidMarketingParam,
		},
		{
			name: "empty history",
			req: &marketing.PointsHistoryRequest{
				UserId:   9999,
				Page:     1,
				PageSize: 10,
			},
			check: func(resp *marketing.PointsHistoryResponse) {
				assert.Equal(s.T(), int64(0), resp.Total)
				assert.Empty(s.T(), resp.Records)
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

			resp, err := s.logic.PointsHistory(tt.req)
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
