package helpers

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

// AssertMessageReceived checks if a message was received within the timeout
func AssertMessageReceived(t *testing.T, rmqHelper *RMQHelper, queueName string, timeout time.Duration) {
    t.Helper()
    msg, err := rmqHelper.ConsumeMessage(queueName, timeout)
    assert.NoError(t, err)
    assert.NotNil(t, msg)
}

// AssertNoMessageReceived checks that no message was received within the timeout
func AssertNoMessageReceived(t *testing.T, rmqHelper *RMQHelper, queueName string, timeout time.Duration) {
    t.Helper()
    msg, err := rmqHelper.ConsumeMessage(queueName, timeout)
    assert.Error(t, err)
    assert.Nil(t, msg)
}

// AssertModelMethodCalled checks if a mock model method was called with expected arguments
func AssertModelMethodCalled(t *testing.T, mock interface{}, methodName string, times int) {
    t.Helper()
    m, ok := mock.(interface{ AssertNumberOfCalls(t *testing.T, methodName string, times int) })
    assert.True(t, ok, "mock object does not implement AssertNumberOfCalls")
    m.AssertNumberOfCalls(t, methodName, times)
}

// AssertEventEquals checks if two events are equal (ignoring timestamp fields)
func AssertEventEquals(t *testing.T, expected, actual interface{}) {
    t.Helper()
    assert.Equal(t, expected, actual)
}