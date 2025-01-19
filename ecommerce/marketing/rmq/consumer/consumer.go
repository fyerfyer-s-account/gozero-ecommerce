package consumer

import (
	"encoding/json"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/consumer/handlers"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

type Consumer struct {
	config            *config.RabbitMQConfig
	conn              *amqp.Connection
	channel           *amqp.Channel
	handlers          map[types.EventType]func(*types.MarketingEvent) error
	mu                sync.RWMutex
	couponsModel      model.CouponsModel
	userCouponsModel  model.UserCouponsModel
	promotionsModel   model.PromotionsModel
	pointsModel       model.UserPointsModel
	pointsRecordModel model.PointsRecordsModel
}

// NewConsumer initializes a new instance of Consumer with provided RabbitMQ configuration and model dependencies.
func NewConsumer(
	cfg *config.RabbitMQConfig,
	couponsModel model.CouponsModel,
	userCouponsModel model.UserCouponsModel,
	promotionsModel model.PromotionsModel,
	pointsModel model.UserPointsModel,
	pointsRecordModel model.PointsRecordsModel,
) (*Consumer, error) {
	c := &Consumer{
		config:            cfg,
		handlers:          make(map[types.EventType]func(*types.MarketingEvent) error),
		couponsModel:      couponsModel,
		userCouponsModel:  userCouponsModel,
		promotionsModel:   promotionsModel,
		pointsModel:       pointsModel,
		pointsRecordModel: pointsRecordModel,
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	c.registerHandlers()
	return c, nil
}

// connect establishes a connection to RabbitMQ and initializes the channel.
func (c *Consumer) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, err := amqp.Dial(c.config.GetDSN())
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		c.config.Exchanges.MarketingEvent.Name,
		c.config.Exchanges.MarketingEvent.Type,
		c.config.Exchanges.MarketingEvent.Durable,
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	c.conn = conn
	c.channel = ch
	return nil
}

// registerHandlers registers the event handlers for different event types.
func (c *Consumer) registerHandlers() {
	// Initialize specific handlers with corresponding models
	couponHandler := handlers.NewCouponHandler(c.couponsModel, c.userCouponsModel)
	promotionHandler := handlers.NewPromotionHandler(c.promotionsModel)
	pointsHandler := handlers.NewPointsHandler(c.pointsModel, c.pointsRecordModel)

	// Map handlers to event types
	c.handlers[types.EventTypeCouponReceived] = couponHandler.HandleCouponReceived
	c.handlers[types.EventTypeCouponUsed] = couponHandler.HandleCouponUsed
	c.handlers[types.EventTypePromotionStarted] = promotionHandler.HandlePromotionStatus
	c.handlers[types.EventTypePromotionEnded] = promotionHandler.HandlePromotionStatus
	c.handlers[types.EventTypePromotionCalculated] = promotionHandler.HandlePromotionCalculated
	c.handlers[types.EventTypePointsAdded] = pointsHandler.HandlePointsTransaction
	c.handlers[types.EventTypePointsUsed] = pointsHandler.HandlePointsTransaction
}

// Start sets up queues, binds them to the exchange, and starts consuming messages.
func (c *Consumer) Start() error {
	for _, q := range []struct {
		name       string
		routingKey string
	}{
		{c.config.Queues.CouponEvent.Name, c.config.Queues.CouponEvent.RoutingKey},
		{c.config.Queues.PromotionEvent.Name, c.config.Queues.PromotionEvent.RoutingKey},
		{c.config.Queues.PointsEvent.Name, c.config.Queues.PointsEvent.RoutingKey},
	} {
		queue, err := c.channel.QueueDeclare(
			q.name,
			true,  // durable
			false, // auto-delete
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return err
		}

		err = c.channel.QueueBind(
			queue.Name,
			q.routingKey,
			c.config.Exchanges.MarketingEvent.Name,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		msgs, err := c.channel.Consume(
			queue.Name,
			"",    // consumer tag
			true,  // auto-ack
			false, // exclusive
			false, // no-local
			false, // no-wait
			nil,   // args
		)
		if err != nil {
			return err
		}

		go c.handle(msgs)
	}

	return nil
}

// handle processes incoming messages by invoking the appropriate handler based on the event type.
func (c *Consumer) handle(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		var event types.MarketingEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		if handler, ok := c.handlers[event.Type]; ok {
			if err := handler(&event); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		} else {
			log.Printf("No handler registered for event type: %v", event.Type)
		}
	}
}

// Close gracefully shuts down the consumer by closing the channel and connection.
func (c *Consumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}
