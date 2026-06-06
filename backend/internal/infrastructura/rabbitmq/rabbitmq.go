package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Client — RabbitMQ connection va channel.
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// New — RabbitMQ'ga ulanadi.
func New(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq.Dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("rabbitmq.Channel: %w", err)
	}
	return &Client{conn: conn, channel: ch}, nil
}

// Close — ulanishni yopadi.
func (c *Client) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

// DeclareQueue — queue'ni e'lon qiladi (idempotent).
func (c *Client) DeclareQueue(name string) error {
	_, err := c.channel.QueueDeclare(
		name,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	return err
}

// Publish — queue'ga message yuboradi.
func (c *Client) Publish(ctx context.Context, queue string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("rabbitmq.Publish marshal: %w", err)
	}
	return c.channel.PublishWithContext(ctx,
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
}

// Consume — queue'dan xabarlarni o'qiydi.
func (c *Client) Consume(queue string) (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		queue,
		"",    // consumer tag
		false, // auto-ack (manual ack for reliability)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
}

const (
	QueueWebhookDelivery = "webhook_delivery"
	QueueEmailSend       = "email_send"
	QueueAutomation      = "automation_event"
)
