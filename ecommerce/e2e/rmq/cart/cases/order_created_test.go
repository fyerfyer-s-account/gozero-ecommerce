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

func TestCartOrderFlow(t *testing.T) {
    suite.Run(t, new(CartTestSuite))
}

func (s *CartTestSuite) TestCartOrderFlow() {
    // Initialize test context
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    testCtx, err := setup.NewTestContext()
    require.NoError(s.T(), err)
    defer testCtx.Close()

    // Load test data
    var event struct {
        Type      string    `json:"type"`
        OrderNo   string    `json:"order_no"`
        UserID    uint64    `json:"user_id"`
        Timestamp time.Time `json:"timestamp"`
        Items     []struct {
            ProductID uint64  `json:"product_id"`
            SkuID    uint64  `json:"sku_id"`
            Quantity int64   `json:"quantity"`
            Price    float64 `json:"price"`
        } `json:"items"`
        TotalAmount float64 `json:"total_amount"`
    }
    require.NoError(s.T(), helpers.LoadFixture("order_created.json", &event))

    // Initialize test cart items
    for _, item := range event.Items {
        cartItem, err := testCtx.DB.PrepareTestCartItem(ctx, event.UserID, item.ProductID, item.SkuID, item.Quantity)
        require.NoError(s.T(), err)
        require.NotNil(s.T(), cartItem)
    }

    // Initialize cart statistics
    stats := &model.CartStatistics{
        UserId:           event.UserID,
        TotalQuantity:    2,
        SelectedQuantity: 2,
        TotalAmount:      event.TotalAmount,
        SelectedAmount:   event.TotalAmount,
    }
    require.NoError(s.T(), testCtx.DB.GetCartStatsModel().Upsert(ctx, stats))

    // Create temporary queue
    queue, err := testCtx.RMQ.GetChannel().QueueDeclare(
        "",    
        false, 
        true,  
        true,  
        false,
        nil,
    )
    require.NoError(s.T(), err)

    // Bind queue to exchange
    err = testCtx.RMQ.GetChannel().QueueBind(
        queue.Name,
        "order.created",
        "order.events",
        false,
        nil,
    )
    require.NoError(s.T(), err)

    // Publish order created event
    err = testCtx.RMQ.PublishEvent("order.events", "order.created", event)
    require.NoError(s.T(), err)

    // Verify message received
    msg := helpers.AssertMessageReceived(s.T(), testCtx.RMQ, queue.Name, 5*time.Second)
    helpers.AssertMessageContent(s.T(), msg, "application/json", "order.created")

    // Wait for processing
    time.Sleep(2 * time.Second)

    // Verify cart is empty
    helpers.AssertCartItemsEmpty(s.T(), testCtx.DB, event.UserID)

    // Verify cart statistics reset
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