package broker

import (
    "context"
    "testing"
    "time"
    
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"
    "github.com/streadway/amqp"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
)

type BrokerTestSuite struct {
    suite.Suite
    broker *AMQPBroker
    ctx    context.Context
}

func TestBrokerSuite(t *testing.T) {
    suite.Run(t, new(BrokerTestSuite))
}

func (s *BrokerTestSuite) SetupSuite() {
    s.ctx = context.Background()
    cfg := config.DefaultConfig()
    s.broker = NewAMQPBroker(cfg)
    err := s.broker.Connect(s.ctx)
    require.NoError(s.T(), err)
}

func (s *BrokerTestSuite) TearDownSuite() {
    s.broker.Close()
}

func (s *BrokerTestSuite) TestExchangeOperations() {
    // Test exchange declaration
    err := s.broker.DeclareExchange("test.exchange", "topic")
    require.NoError(s.T(), err)

    // Test declaring same exchange again should not error
    err = s.broker.DeclareExchange("test.exchange", "topic")
    require.NoError(s.T(), err)
}

func (s *BrokerTestSuite) TestQueueOperations() {
    // Test queue declaration
    queue, err := s.broker.DeclareQueue("test.queue")
    require.NoError(s.T(), err)
    require.Equal(s.T(), "test.queue", queue.Name)

    // Test queue binding
    err = s.broker.BindQueue("test.queue", "test.exchange", "test.*")
    require.NoError(s.T(), err)
}

func (s *BrokerTestSuite) TestPublishAndConsume() {
    // Setup
    exchangeName := "test.pub.exchange"
    queueName := "test.pub.queue"
    routingKey := "test.message"
    
    err := s.broker.DeclareExchange(exchangeName, "topic")
    require.NoError(s.T(), err)
    
    _, err = s.broker.DeclareQueue(queueName)
    require.NoError(s.T(), err)
    
    err = s.broker.BindQueue(queueName, exchangeName, routingKey)
    require.NoError(s.T(), err)

    // Start consumer
    msgs, err := s.broker.Consume(queueName, true)
    require.NoError(s.T(), err)

    // Publish message
    message := "test message"
    err = s.broker.Publish(exchangeName, routingKey, amqp.Publishing{
        ContentType: "text/plain",
        Body:        []byte(message),
    })
    require.NoError(s.T(), err)

    // Verify message received
    select {
    case msg := <-msgs:
        require.Equal(s.T(), message, string(msg.Body))
    case <-time.After(5 * time.Second):
        s.T().Fatal("timeout waiting for message")
    }
}

func (s *BrokerTestSuite) TestChannelOperations() {
    // Test getting new channel
    ch1, err := s.broker.Channel()
    require.NoError(s.T(), err)
    require.NotNil(s.T(), ch1)
    defer ch1.Close()

    // Test getting another channel
    ch2, err := s.broker.Channel()
    require.NoError(s.T(), err)
    require.NotNil(s.T(), ch2)
    defer ch2.Close()

    // Verify channels are different
    require.NotEqual(s.T(), ch1, ch2)

    // Test channel operations
    err = ch1.ExchangeDeclare(
        "test.channel.exchange",
        "topic",
        true,
        false,
        false,
        false,
        nil,
    )
    require.NoError(s.T(), err)

    // Test publishing through channel
    err = ch1.Publish(
        "test.channel.exchange",
        "test.key",
        false,
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:       []byte("test message"),
        },
    )
    require.NoError(s.T(), err)
}

func (s *BrokerTestSuite) TestDisconnect() {
    // Setup new broker instance
    cfg := config.DefaultConfig()
    broker := NewAMQPBroker(cfg)
    
    // Connect
    err := broker.Connect(s.ctx)
    require.NoError(s.T(), err)
    require.True(s.T(), broker.IsConnected())
    
    // Test disconnect
    err = broker.Disconnect()
    require.NoError(s.T(), err)
    require.False(s.T(), broker.IsConnected())
    
    // Verify operations fail after disconnect
    _, err = broker.DeclareQueue("test.disconnect.queue")
    require.Error(s.T(), err)
    
    // Test multiple disconnects don't error
    err = broker.Disconnect()
    require.NoError(s.T(), err)
    
    // Test can reconnect after disconnect
    err = broker.Connect(s.ctx)
    require.NoError(s.T(), err)
    require.True(s.T(), broker.IsConnected())
}