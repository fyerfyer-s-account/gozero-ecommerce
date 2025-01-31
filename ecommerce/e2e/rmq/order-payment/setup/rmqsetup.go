package setup

import (
    "context"
    "fmt"

    "github.com/streadway/amqp"
)

type RMQTestSetup struct {
    Conn    *amqp.Connection
    Channel *amqp.Channel
    Config  *TestConfig
}

func NewRMQTestSetup(config *TestConfig) (*RMQTestSetup, error) {
    // Connect to RabbitMQ
    conn, err := amqp.Dial(config.RabbitMQURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    // Create channel
    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("failed to open channel: %w", err)
    }

    return &RMQTestSetup{
        Conn:    conn,
        Channel: ch,
        Config:  config,
    }, nil
}

func (s *RMQTestSetup) SetupTestEnvironment(ctx context.Context) error {
    // Declare exchanges
    if err := s.declareExchanges(); err != nil {
        return err
    }

    // Declare queues and bindings
    if err := s.declareQueuesAndBindings(); err != nil {
        return err
    }

    return nil
}

func (s *RMQTestSetup) declareExchanges() error {
    // Declare order events exchange
    err := s.Channel.ExchangeDeclare(
        "order.events",  // name
        "topic",         // type
        true,           // durable
        false,          // auto-deleted
        false,          // internal
        false,          // no-wait
        nil,            // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare order exchange: %w", err)
    }

    // Declare payment events exchange
    err = s.Channel.ExchangeDeclare(
        "payment.events", // name
        "topic",         // type
        true,           // durable
        false,          // auto-deleted
        false,          // internal
        false,          // no-wait
        nil,            // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare payment exchange: %w", err)
    }

    return nil
}

func (s *RMQTestSetup) declareQueuesAndBindings() error {
    // Declare and bind order.payment.success queue
    if _, err := s.Channel.QueueDeclare(
        "order.payment.success", // name
        true,                   // durable
        false,                  // delete when unused
        false,                  // exclusive
        false,                  // no-wait
        nil,                    // arguments
    ); err != nil {
        return fmt.Errorf("failed to declare order.payment.success queue: %w", err)
    }

    err := s.Channel.QueueBind(
        "order.payment.success", // queue name
        "payment.success",       // routing key
        "payment.events",        // exchange
        false,
        nil,
    )
    if err != nil {
        return fmt.Errorf("failed to bind order.payment.success queue: %w", err)
    }

    // Declare and bind order.payment.failed queue
    if _, err := s.Channel.QueueDeclare(
        "order.payment.failed", // name
        true,                  // durable
        false,                 // delete when unused
        false,                 // exclusive
        false,                 // no-wait
        nil,                   // arguments
    ); err != nil {
        return fmt.Errorf("failed to declare order.payment.failed queue: %w", err)
    }

    err = s.Channel.QueueBind(
        "order.payment.failed", // queue name
        "payment.failed",      // routing key
        "payment.events",      // exchange
        false,
        nil,
    )
    if err != nil {
        return fmt.Errorf("failed to bind order.payment.failed queue: %w", err)
    }

    return nil
}

func (s *RMQTestSetup) CleanupTestQueues() error {
    queues := []string{
        "order.payment.success",
        "order.payment.failed",
    }

    for _, queue := range queues {
        _, err := s.Channel.QueuePurge(queue, false)
        if err != nil {
            return fmt.Errorf("failed to purge queue %s: %w", queue, err)
        }
    }

    return nil
}

func (s *RMQTestSetup) Close() error {
    if s.Channel != nil {
        if err := s.Channel.Close(); err != nil {
            return fmt.Errorf("failed to close channel: %w", err)
        }
    }
    
    if s.Conn != nil {
        if err := s.Conn.Close(); err != nil {
            return fmt.Errorf("failed to close connection: %w", err)
        }
    }
    
    return nil
}