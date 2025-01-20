package logic

import (
	"context"
	"database/sql"
	"encoding/json"
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

type CalculatePromotionTestSuite struct {
	suite.Suite
	ctx   *svc.ServiceContext
	logic *CalculatePromotionLogic
}

func TestCalculatePromotionSuite(t *testing.T) {
	suite.Run(t, new(CalculatePromotionTestSuite))
}

func (s *CalculatePromotionTestSuite) SetupSuite() {
	configFile := flag.String("f", "../../etc/marketing.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
	s.logic = NewCalculatePromotionLogic(context.Background(), s.ctx)
}

func (s *CalculatePromotionTestSuite) TearDownTest() {
	s.cleanData()
}

func (s *CalculatePromotionTestSuite) cleanData() {
	_ = s.ctx.PromotionsModel.Delete(context.Background(), 1001)
	_ = s.ctx.PromotionsModel.Delete(context.Background(), 1002)
}

func (s *CalculatePromotionTestSuite) TestCalculatePromotion() {
	tests := []struct {
		name    string
		setup   func() error
		req     *marketing.CalculatePromotionRequest
		want    *marketing.CalculatePromotionResponse
		wantErr error
	}{
		{
			name: "fixed amount discount",
			setup: func() error {
				rule := PromotionRule{
					MinAmount:      100,
					DiscountAmount: 20,
				}
				ruleJson, _ := json.Marshal(rule)

				now := time.Now()
				_, err := s.ctx.PromotionsModel.Insert(context.Background(), &model.Promotions{
					Id:        1001,
					Name:      "20元优惠",
					Type:      1,
					Rules:     string(ruleJson),
					Status:    1,
					StartTime: sql.NullTime{Time: now.Add(-time.Hour), Valid: true},
					EndTime:   sql.NullTime{Time: now.Add(time.Hour), Valid: true},
				})
				return err
			},
			req: &marketing.CalculatePromotionRequest{
				Items: []*marketing.OrderItem{
					{ProductId: 1, Price: 150, Quantity: 1},
				},
			},
			want: &marketing.CalculatePromotionResponse{
				OriginalAmount: 150,
				DiscountAmount: 20,
				FinalAmount:    130,
			},
		},
		{
			name: "percentage discount with max limit",
			setup: func() error {
				rule := PromotionRule{
					MinAmount:         200,
					DiscountRate:      0.2,
					MaxDiscountAmount: 50,
				}
				ruleJson, _ := json.Marshal(rule)

				now := time.Now()
				_, err := s.ctx.PromotionsModel.Insert(context.Background(), &model.Promotions{
					Id:        1002,
					Name:      "8折优惠",
					Type:      2,
					Rules:     string(ruleJson),
					Status:    1,
					StartTime: sql.NullTime{Time: now.Add(-time.Hour), Valid: true},
					EndTime:   sql.NullTime{Time: now.Add(time.Hour), Valid: true},
				})
				return err
			},
			req: &marketing.CalculatePromotionRequest{
				Items: []*marketing.OrderItem{
					{ProductId: 1, Price: 300, Quantity: 1},
				},
			},
			want: &marketing.CalculatePromotionResponse{
				OriginalAmount: 300,
				DiscountAmount: 50,
				FinalAmount:    250,
			},
		},
		{
			name: "empty items",
			req: &marketing.CalculatePromotionRequest{
				Items: []*marketing.OrderItem{},
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

			got, err := s.logic.CalculatePromotion(tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}

			assert.NoError(t, err)
			if tt.want != nil {
				assert.Equal(t, tt.want.OriginalAmount, got.OriginalAmount)
				assert.Equal(t, tt.want.DiscountAmount, got.DiscountAmount)
				assert.Equal(t, tt.want.FinalAmount, got.FinalAmount)
			}
		})
	}
}
