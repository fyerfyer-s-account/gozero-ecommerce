package cases

import (
    "context"
    "testing"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/e2e/rmq/cart/setup"
    "github.com/stretchr/testify/suite"
)

type CartTestSuite struct {
    suite.Suite
    ctx     context.Context
    cancel  context.CancelFunc
    testCtx *setup.TestContext
}

func (s *CartTestSuite) SetupSuite() {
    s.ctx, s.cancel = context.WithTimeout(context.Background(), 30*time.Second)
    testCtx, err := setup.NewTestContext()
    s.Require().NoError(err)
    s.testCtx = testCtx
}

func (s *CartTestSuite) TearDownSuite() {
    s.cancel()
    if s.testCtx != nil {
        s.Require().NoError(s.testCtx.Close())
    }
}

func (s *CartTestSuite) SetupTest() {
    s.Require().NoError(s.testCtx.DB.CleanTestData(s.ctx))
}

func (s *CartTestSuite) TearDownTest() {
    s.Require().NoError(s.testCtx.DB.CleanTestData(s.ctx))
}

func TestCartSuite(t *testing.T) {
    suite.Run(t, new(CartTestSuite))
}