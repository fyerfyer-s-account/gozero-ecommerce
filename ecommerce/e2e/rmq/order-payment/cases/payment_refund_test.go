package cases

import (
	// "context"
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

func TestPaymentRefundFlow(t *testing.T) {
    // Setup
    orderPaymentsModel := &mocks.OrderPaymentsModel{}
    ordersModel := &mocks.OrdersModel{}
    orderRefundsModel := &mocks.OrderRefundsModel{}

    rmqHelper, err := helpers.NewRMQHelper()
    assert.NoError(t, err)
    defer rmqHelper.Close()

    // Load test event from fixture
    var event types.PaymentRefundEvent
    err = helpers.LoadFixture("payment_refund.json", &event)
    assert.NoError(t, err)

    // Setup expectations
    ordersModel.On("FindByOrderNo", mock.Anything, event.OrderNo).Return(&model.Orders{
        Id:     1,
        Status: 4, // Completed
    }, nil)
    ordersModel.On("UpdateStatus", mock.Anything, uint64(1), int64(6)).Return(nil)

    orderPaymentsModel.On("FindOneByPaymentNo", mock.Anything, event.PaymentNo).Return(&model.OrderPayments{
        Id:     1,
        Status: 1, // Paid
    }, nil)
    orderPaymentsModel.On("UpdateStatus", mock.Anything, event.PaymentNo, 2, mock.Anything).Return(nil)

    orderRefundsModel.On("FindOneByRefundNo", mock.Anything, event.RefundNo).Return(&model.OrderRefunds{
        Id:     1,
        Status: 2, // Processing
    }, nil)
    orderRefundsModel.On("UpdateStatus", mock.Anything, event.RefundNo, 3, "Refund completed").Return(nil)

    eventJSON, err := json.Marshal(event)
    assert.NoError(t, err)

    // Publish event
    err = rmqHelper.PublishMessage("payment.events", "payment.success", eventJSON)
    assert.NoError(t, err)

    // Wait and verify order service received and processed the message
    helpers.AssertMessageReceived(t, rmqHelper, "order.payment.refund", 5*time.Second)

    // Verify model calls
    ordersModel.AssertExpectations(t)
    orderPaymentsModel.AssertExpectations(t)
    orderRefundsModel.AssertExpectations(t)
}