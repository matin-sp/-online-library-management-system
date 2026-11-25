package queue

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    amqp091 "github.com/rabbitmq/amqp091-go"
)


type ReturnConfirmationMessage struct {
    LoanID int64  `json:"loan_id"`
    UserID int64  `json:"user_id"`
    Email  string `json:"email"`
}


type Publisher interface {
    PublishReturnConfirmation(msg ReturnConfirmationMessage) error
    Close() error
}


type RabbitMQPublisher struct {
    conn      *amqp091.Connection
    channel   *amqp091.Channel
    queueName string
}

// NewRabbitMQPublisher dials the broker, opens a channel and makes sure
// the queue exists. Fail-fast: if RabbitMQ is down, the app won't start.
func NewRabbitMQPublisher(url string, queueName string) (*RabbitMQPublisher, error) {

    conn, err := amqp091.Dial(url)
    if err != nil {
        return nil, fmt.Errorf("connect to RabbitMQ failed: %w", err)
    }


    channel, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("open RabbitMQ channel failed: %w", err)
    }


    _, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("declare queue %s failed: %w", queueName, err)
    }

    return &RabbitMQPublisher{
        conn:      conn,
        channel:   channel,
        queueName: queueName,
    }, nil
}


func (p *RabbitMQPublisher) PublishReturnConfirmation(msg ReturnConfirmationMessage) error {
    body, err := json.Marshal(msg)
    if err != nil {
        return fmt.Errorf("marshal message to JSON failed: %w", err)
    }


    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()


    err = p.channel.PublishWithContext(ctx,
        "",
        p.queueName,
        false, 
        false, 
        amqp091.Publishing{
            ContentType: "application/json",
            Body:        body,
            // Persistent -> RabbitMQ writes the message to disk.
            // Together with the durable queue, the message survives
            // a broker restart.
            DeliveryMode: amqp091.Persistent,
        },
    )
    if err != nil {
        return fmt.Errorf("publish message to queue failed: %w", err)
    }
    return nil
}


func (p *RabbitMQPublisher) Close() error {
    if p.channel != nil {
        p.channel.Close()
    }
    if p.conn != nil {
        return p.conn.Close()
    }
    return nil
}