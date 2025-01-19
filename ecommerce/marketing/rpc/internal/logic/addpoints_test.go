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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type AddPointsTestSuite struct {
	suite.Suite
	ctx   *svc.ServiceContext
	logic *AddPointsLogic
}

func TestAddPointsSuite(t *testing.T) {
	suite.Run(t, new(AddPointsTestSuite))
}

func (s *AddPointsTestSuite) SetupSuite() {
	// Load configuration file
	configFile := flag.String("f", "../../etc/marketing.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
	s.logic = NewAddPointsLogic(context.Background(), s.ctx)
}

func (s *AddPointsTestSuite) SetupTest() {
	// Clean up data before each test
	s.cleanData()
}

func (s *AddPointsTestSuite) TearDownTest() {
	// Clean up data after each test
	s.cleanData()
}

func (s *AddPointsTestSuite) cleanData() {
	// Ensure cleanup runs in a transaction
	ctx := context.Background()
	_ = s.ctx.UserPointsModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
		_ = s.ctx.UserPointsModel.Delete(ctx, 1001)
		_ = s.ctx.UserPointsModel.Delete(ctx, 1002)
		return nil
	})
}

func (s *AddPointsTestSuite) TestAddPoints() {
	tests := []struct {
		name    string
		setup   func() error
		req     *marketing.AddPointsRequest
		wantErr error
		check   func(*marketing.AddPointsResponse)
	}{
		{
			name: "add points to new user",
			req: &marketing.AddPointsRequest{
				UserId:  1001,
				Points:  100,
				Source:  "test",
				Remark:  "test points",
			},
			wantErr: nil,
			check: func(resp *marketing.AddPointsResponse) {
				// Check that the response is not nil
				if resp == nil {
					s.T().Error("Response should not be nil")
					return
				}
				assert.True(s.T(), resp.Success)
				assert.Equal(s.T(), int64(100), resp.CurrentPoints)

				// Verify database balance for user 1001
				points, err := s.ctx.UserPointsModel.GetBalance(context.Background(), 1001)
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), int64(100), points)
			},
		},
		{
			name: "add points to existing user",
			setup: func() error {
				// Setup test by inserting initial points for user 1002
				return s.ctx.UserPointsModel.Trans(context.Background(), func(ctx context.Context, session sqlx.Session) error {
					_, err := s.ctx.UserPointsModel.Insert(ctx, &model.UserPoints{
						UserId:      1002,
						Points:      100,
						TotalPoints: 100,
						UsedPoints:  0,
						UpdatedAt:   time.Now(),
					})
					return err
				})
			},
			req: &marketing.AddPointsRequest{
				UserId:  1002,
				Points:  50,
				Source:  "test",
				Remark:  "add more points",
			},
			wantErr: nil,
			check: func(resp *marketing.AddPointsResponse) {
				// Check that the response is correct
				assert.True(s.T(), resp.Success)
				assert.Equal(s.T(), int64(150), resp.CurrentPoints)

				// Verify database balance for user 1002
				points, err := s.ctx.UserPointsModel.GetBalance(context.Background(), 1002)
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), int64(150), points)
			},
		},
		{
			name: "invalid points amount",
			req: &marketing.AddPointsRequest{
				UserId:  1001,
				Points:  -100,
				Source:  "test",
				Remark:  "negative points",
			},
			wantErr: zeroerr.ErrInvalidPointsAmount,
		},
		{
			name: "exceed points limit",
			req: &marketing.AddPointsRequest{
				UserId:  1001,
				Points:  1000001,
				Source:  "test",
				Remark:  "too many points",
			},
			wantErr: zeroerr.ErrExceedPointsLimit,
		},
	}

	// Loop over each test case and run them
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Clean data before running each test
			s.cleanData()
			if tt.setup != nil {
				err := tt.setup()
				assert.NoError(t, err)
			}

			// Call AddPoints logic
			resp, err := s.logic.AddPoints(tt.req)

			// Check for expected errors
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}

			// If no error, check the response and verify the points
			assert.NoError(t, err)
			if tt.check != nil {
				tt.check(resp)
			}
		})
	}
}
