package cases

import (
	// "context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/e2e/rmq/order-payment/helpers"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/e2e/rmq/order-payment/mocks"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPaymentFailedFlow(t *testing.T) {
    // Setup
    orderPaymentsModel := &mocks.OrderPaymentsModel{}
    ordersModel := &mocks.OrdersModel{}

    rmqHelper, err := helpers.NewRMQHelper()
    assert.NoError(t, err)
    defer rmqHelper.Close()

    // Load test event from fixture
    var event types.PaymentFailedEvent
    err = helpers.LoadFixture("payment_failed.json", &event)
    assert.NoError(t, err)

    // Setup expectations
    ordersModel.On("FindByOrderNo", mock.Anything, event.OrderNo).Return(&model.Orders{
        Id:     1,
        Status: 1, // Pending payment
    }, nil)
    ordersModel.On("UpdateStatus", mock.Anything, uint64(1), int64(5)).Return(nil)

    orderPaymentsModel.On("FindOneByPaymentNo", mock.Anything, event.PaymentNo).Return(&model.OrderPayments{
        Id:     1,
        Status: 1,
        PayTime: sql.NullTime{Time: time.Now(), Valid: true},
    }, nil)
    orderPaymentsModel.On("UpdateStatus", mock.Anything, event.PaymentNo, 0, mock.Anything).Return(nil)

    eventJSON, err := json.Marshal(event)
    assert.NoError(t, err)

    // Publish event
    err = rmqHelper.PublishMessage("payment.events", "payment.failed", eventJSON)
    assert.NoError(t, err)

    // Wait and verify order service received and processed the message
    helpers.AssertMessageReceived(t, rmqHelper, "order.payment.failed", 5*time.Second)

    // Verify model calls
    ordersModel.AssertExpectations(t)
    orderPaymentsModel.AssertExpectations(t)
}