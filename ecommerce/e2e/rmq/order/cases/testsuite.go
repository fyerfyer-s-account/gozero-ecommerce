package cases

import (
    "context"
    "testing"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/e2e/rmq/order/setup"
    "github.com/stretchr/testify/suite"
)

type PaymentTestSuite struct {
    suite.Suite
    ctx     context.Context
    cancel  context.CancelFunc
    testCtx *setup.TestContext
}

func (s *PaymentTestSuite) SetupSuite() {
    s.ctx, s.cancel = context.WithTimeout(context.Background(), 30*time.Second)
    testCtx, err := setup.NewTestContext()
    s.Require().NoError(err)
    s.testCtx = testCtx
}

func (s *PaymentTestSuite) TearDownSuite() {
    s.cancel()
    if s.testCtx != nil {
        s.Require().NoError(s.testCtx.Close())
    }
}

func (s *PaymentTestSuite) SetupTest() {
    // Clean test data before each test
    s.Require().NoError(s.testCtx.DB.CleanTestData(s.ctx))
}

func (s *PaymentTestSuite) TearDownTest() {
    // Clean test data after each test
    s.Require().NoError(s.testCtx.DB.CleanTestData(s.ctx))
}

func TestPaymentSuite(t *testing.T) {
    suite.Run(t, new(PaymentTestSuite))
}