package cases

import (
	"context"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/e2e/rmq/cart/helpers"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/e2e/rmq/cart/setup"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestCartPaymentFlow(t *testing.T) {
    suite.Run(t, new(CartTestSuite))
}

func (s *CartTestSuite) TestCartPaymentFlow() {
    // Initialize test context
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    testCtx, err := setup.NewTestContext()
    require.NoError(s.T(), err)
    defer testCtx.Close()

    // Load test data
    var event struct {
        Type          string    `json:"type"`
        OrderNo       string    `json:"order_no"`
        PaymentNo     string    `json:"payment_no"`
        UserID        uint64    `json:"user_id"`
        Timestamp     time.Time `json:"timestamp"`
        Amount        float64   `json:"amount"`
        PaymentMethod int       `json:"payment_method"`
        PaidTime      time.Time `json:"paid_time"`
    }
    require.NoError(s.T(), helpers.LoadFixture("payment_success.json", &event))

    // Initialize test cart items
    cartItem, err := testCtx.DB.PrepareTestCartItem(ctx, event.UserID, 100, 1001, 2)
    require.NoError(s.T(), err)
    require.NotNil(s.T(), cartItem)

    // Initialize cart statistics
    stats := &model.CartStatistics{
        UserId:           event.UserID,
        TotalQuantity:    2,
        SelectedQuantity: 2,
        TotalAmount:      event.Amount,
        SelectedAmount:   event.Amount,
    }
    require.NoError(s.T(), testCtx.DB.GetCartStatsModel().Upsert(ctx, stats))

    // Publish payment success event
    s.T().Log("Publishing payment success event")
    err = testCtx.RMQ.PublishEvent("payment.events", "payment.success", event)
    require.NoError(s.T(), err)

    // Wait for cart service to process the message
    time.Sleep(2 * time.Second)

    // Verify cart is cleared
    s.T().Log("Verifying cart data")
    helpers.AssertCartItemsEmpty(s.T(), testCtx.DB, event.UserID)

    // Verify cart statistics are reset
    helpers.AssertCartStatistics(s.T(), testCtx.DB, event.UserID, struct {
        TotalQuantity    int64
        SelectedQuantity int64
        TotalAmount      float64
        SelectedAmount   float64
    }{
        TotalQuantity:    0,
        SelectedQuantity: 0,
        TotalAmount:      0,
        SelectedAmount:   0,
    })
}