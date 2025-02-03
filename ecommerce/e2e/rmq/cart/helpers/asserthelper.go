package helpers

import (
    "context"
    "testing"
    "time"

    "github.com/streadway/amqp"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// RMQAssertions defines interface for RabbitMQ related assertions
type RMQAssertions interface {
    ConsumeMessage(queueName string, timeout time.Duration) (*amqp.Delivery, error)
}

// DBAssertions defines interface for database related assertions
type DBAssertions interface {
    GetCartItems(ctx context.Context, userId uint64) (interface{}, error)
    GetCartStatistics(ctx context.Context, userId uint64) (interface{}, error)
}

// AssertMessageReceived checks if a message was received on the specified queue
func AssertMessageReceived(t *testing.T, rmq RMQAssertions, queueName string, timeout time.Duration) *amqp.Delivery {
    t.Helper()
    msg, err := rmq.ConsumeMessage(queueName, timeout)
    require.NoError(t, err, "failed to consume message from queue %s", queueName)
    require.NotNil(t, msg, "no message received from queue %s within timeout", queueName)
    return msg
}

// AssertNoMessageReceived verifies that no message is received within timeout
func AssertNoMessageReceived(t *testing.T, rmq RMQAssertions, queueName string, timeout time.Duration) {
    t.Helper()
    msg, err := rmq.ConsumeMessage(queueName, timeout)
    assert.Error(t, err, "unexpected message received from queue %s", queueName)
    assert.Nil(t, msg, "expected no message but got one from queue %s", queueName)
}

// AssertCartItemsEmpty verifies that the user's cart is empty
func AssertCartItemsEmpty(t *testing.T, db *DBHelper, userId uint64) {
    t.Helper()
    ctx := context.Background()
    
    items, err := db.GetCartItemsModel().FindByUserId(ctx, userId)
    require.NoError(t, err, "failed to get cart items for user %d", userId)
    assert.Empty(t, items, "cart should be empty for user %d", userId)
}

// AssertCartStatistics verifies the cart statistics match expected values
func AssertCartStatistics(t *testing.T, db *DBHelper, userId uint64, expectedStats struct {
    TotalQuantity    int64
    SelectedQuantity int64
    TotalAmount      float64
    SelectedAmount   float64
}) {
    t.Helper()
    ctx := context.Background()
    
    stats, err := db.GetCartStatsModel().FindOne(ctx, userId)
    require.NoError(t, err, "failed to get cart statistics for user %d", userId)
    require.NotNil(t, stats, "cart statistics not found for user %d", userId)
    
    assert.Equal(t, expectedStats.TotalQuantity, stats.TotalQuantity,
        "total quantity mismatch for user %d", userId)
    assert.Equal(t, expectedStats.SelectedQuantity, stats.SelectedQuantity,
        "selected quantity mismatch for user %d", userId)
    assert.Equal(t, expectedStats.TotalAmount, stats.TotalAmount,
        "total amount mismatch for user %d", userId)
    assert.Equal(t, expectedStats.SelectedAmount, stats.SelectedAmount,
        "selected amount mismatch for user %d", userId)
}

// AssertMessageContent verifies the message content matches expected data
func AssertMessageContent(t *testing.T, msg *amqp.Delivery, expectedType string, expectedRoutingKey string) {
    t.Helper()
    assert.Equal(t, "application/json", msg.ContentType, "unexpected message content type")
    assert.Equal(t, expectedRoutingKey, msg.RoutingKey, "unexpected routing key")
    assert.NotEmpty(t, msg.Body, "message body should not be empty")
}

// AssertCartItemSelected verifies if specific cart items are selected/unselected
func AssertCartItemSelected(t *testing.T, db *DBHelper, userId uint64, skuId uint64, expectedSelected int64) {
    t.Helper()
    ctx := context.Background()
    
    item, err := db.GetCartItemsModel().FindOneByUserIdSkuId(ctx, userId, skuId)
    require.NoError(t, err, "failed to get cart item for user %d and sku %d", userId, skuId)
    require.NotNil(t, item, "cart item not found for user %d and sku %d", userId, skuId)
    
    assert.Equal(t, expectedSelected, item.Selected,
        "cart item selection mismatch for user %d and sku %d", userId, skuId)
}

// AssertDatabaseState verifies all relevant database states in one call
func AssertDatabaseState(t *testing.T, db *DBHelper, userId uint64, expectedItems int, expectedStats struct {
    TotalQuantity    int64
    SelectedQuantity int64
    TotalAmount      float64
    SelectedAmount   float64
}) {
    t.Helper()
    ctx := context.Background()
    
    // Check cart items count
    items, err := db.GetCartItemsModel().FindByUserId(ctx, userId)
    require.NoError(t, err, "failed to get cart items")
    assert.Len(t, items, expectedItems, "unexpected number of cart items")
    
    // Check cart statistics
    stats, err := db.GetCartStatsModel().FindOne(ctx, userId)
    require.NoError(t, err, "failed to get cart statistics")
    if expectedItems == 0 {
        assert.Nil(t, stats, "cart statistics should be nil when cart is empty")
        return
    }
    
    require.NotNil(t, stats, "cart statistics should not be nil")
    assert.Equal(t, expectedStats.TotalQuantity, stats.TotalQuantity, "total quantity mismatch")
    assert.Equal(t, expectedStats.SelectedQuantity, stats.SelectedQuantity, "selected quantity mismatch")
    assert.Equal(t, expectedStats.TotalAmount, stats.TotalAmount, "total amount mismatch")
    assert.Equal(t, expectedStats.SelectedAmount, stats.SelectedAmount, "selected amount mismatch")
}