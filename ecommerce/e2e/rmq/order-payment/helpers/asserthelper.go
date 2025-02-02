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
    GetOrderByNo(ctx context.Context, orderNo string) (interface{}, error)
    GetPaymentByNo(ctx context.Context, paymentNo string) (interface{}, error)
    GetRefundByNo(ctx context.Context, refundNo string) (interface{}, error)
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

// AssertOrderStatus verifies the order status
func AssertOrderStatus(t *testing.T, db *DBHelper, orderNo string, expectedStatus int64) {
    t.Helper()
    ctx := context.Background()
    
    order, err := db.GetOrdersModel().FindByOrderNo(ctx, orderNo)
    require.NoError(t, err, "failed to get order by order number %s", orderNo)
    require.NotNil(t, order, "order not found with order number %s", orderNo)
    
    assert.Equal(t, expectedStatus, order.Status, 
        "order status mismatch for order %s: expected %d, got %d", 
        orderNo, expectedStatus, order.Status)
}

// AssertPaymentStatus verifies the payment status
func AssertPaymentStatus(t *testing.T, db *DBHelper, paymentNo string, expectedStatus int64) {
    t.Helper()
    ctx := context.Background()
    
    payment, err := db.GetPaymentsModel().FindOneByPaymentNo(ctx, paymentNo)
    require.NoError(t, err, "failed to get payment by payment number %s", paymentNo)
    require.NotNil(t, payment, "payment not found with payment number %s", paymentNo)
    
    assert.Equal(t, expectedStatus, payment.Status,
        "payment status mismatch for payment %s: expected %d, got %d",
        paymentNo, expectedStatus, payment.Status)
}

// AssertRefundStatus verifies the refund status
func AssertRefundStatus(t *testing.T, db *DBHelper, refundNo string, expectedStatus int64) {
    t.Helper()
    ctx := context.Background()
    
    refund, err := db.GetRefundsModel().FindOneByRefundNo(ctx, refundNo)
    require.NoError(t, err, "failed to get refund by refund number %s", refundNo)
    require.NotNil(t, refund, "refund not found with refund number %s", refundNo)
    
    assert.Equal(t, expectedStatus, refund.Status,
        "refund status mismatch for refund %s: expected %d, got %d",
        refundNo, expectedStatus, refund.Status)
}

// AssertMessageContent verifies the message content matches expected data
func AssertMessageContent(t *testing.T, msg *amqp.Delivery, expectedType string, expectedRoutingKey string) {
    t.Helper()
    assert.Equal(t, "application/json", msg.ContentType, "unexpected message content type")
    assert.Equal(t, expectedRoutingKey, msg.RoutingKey, "unexpected routing key")
    assert.NotEmpty(t, msg.Body, "message body should not be empty")
}

// AssertDatabaseState verifies multiple database conditions in one call
func AssertDatabaseState(t *testing.T, db *DBHelper, checks map[string]int64) {
    t.Helper()
    
    for id, expectedStatus := range checks {
        switch {
        case id[0:5] == "TEST_": // Order number format
            // Check if it's an order number
            if len(id) >= 8 && id[5:8] == "ORD" {
                AssertOrderStatus(t, db, id, expectedStatus)
            // Check if it's a payment number    
            } else if len(id) >= 8 && id[5:8] == "PAY" {
                AssertPaymentStatus(t, db, id, expectedStatus)
            // Check if it's a refund number    
            } else if len(id) >= 8 && id[5:8] == "REF" {
                AssertRefundStatus(t, db, id, expectedStatus)
            } else {
                t.Errorf("invalid ID format: %s", id)
            }
        default:
            t.Errorf("ID must start with TEST_: %s", id)
        }
    }
}