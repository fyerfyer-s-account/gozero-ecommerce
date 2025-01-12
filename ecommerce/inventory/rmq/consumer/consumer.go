package consumer

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/consumer/handlers"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventoryclient"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/messageservice"
	"github.com/streadway/amqp"
)

type Consumer struct {
	config       *config.RabbitMQConfig
	conn         *amqp.Connection
	channel      *amqp.Channel
	handlers     map[types.EventType]func(*types.InventoryEvent) error
	mu           sync.RWMutex
	inventoryRpc inventoryclient.Inventory    
	messageRpc   messageservice.MessageService
}

func NewConsumer(
	cfg *config.RabbitMQConfig,
	inventoryRpc inventoryclient.Inventory,
	messageRpc messageservice.MessageService,
) (*Consumer, error) {
	c := &Consumer{
		config:       cfg,
		handlers:     make(map[types.EventType]func(*types.InventoryEvent) error),
		inventoryRpc: inventoryRpc,
		messageRpc:   messageRpc,
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	c.registerHandlers()
	return c, nil
}

func (c *Consumer) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, err := amqp.Dial(c.config.GetDSN())
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	err = ch.ExchangeDeclare(
		c.config.Exchanges.InventoryEvent.Name,
		c.config.Exchanges.InventoryEvent.Type,
		c.config.Exchanges.InventoryEvent.Durable,
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return err
	}

	c.conn = conn
	c.channel = ch
	return nil
}

func (c *Consumer) registerHandlers() {
	stockUpdateHandler := handlers.NewStockUpdateHandler(c.inventoryRpc)
	stockAlertHandler := handlers.NewStockAlertHandler(c.inventoryRpc, c.messageRpc)
	stockLockHandler := handlers.NewStockLockHandler(c.inventoryRpc)

	c.handlers[types.EventTypeStockUpdated] = stockUpdateHandler.Handle
	c.handlers[types.EventTypeStockAlert] = stockAlertHandler.Handle
	c.handlers[types.EventTypeStockLocked] = stockLockHandler.Handle
	c.handlers[types.EventTypeStockUnlocked] = stockLockHandler.Handle
}

func (c *Consumer) Start() error {
	queues := []struct {
		name       string
		routingKey string
	}{
		{
			name:       c.config.Queues.StockUpdate.Name,
			routingKey: c.config.Queues.StockUpdate.RoutingKey,
		},
		{
			name:       c.config.Queues.StockAlert.Name,
			routingKey: c.config.Queues.StockAlert.RoutingKey,
		},
		{
			name:       c.config.Queues.StockLock.Name,
			routingKey: c.config.Queues.StockLock.RoutingKey,
		},
	}

	for _, q := range queues {
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
			c.config.Exchanges.InventoryEvent.Name,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		msgs, err := c.channel.Consume(
			queue.Name,
			"",    // consumer
			false, // auto-ack
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

func (c *Consumer) handle(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		var event types.InventoryEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			msg.Nack(false, true) // requeue message
			continue
		}

		if handler, ok := c.handlers[event.Type]; ok {
			if err := handler(&event); err != nil {
				log.Printf("Error handling message: %v", err)
				msg.Nack(false, true)
				continue
			}
			msg.Ack(false)
		} else {
			log.Printf("No handler for event type: %s", event.Type)
			msg.Nack(false, false) // don't requeue unknown events
		}
	}
}

func (c *Consumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}

	return nil
}
