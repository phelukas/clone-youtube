package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

func newConnection(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to Rabbitmq: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open to Rabbitmq: %v", err)
	}

	return conn, channel, nil
}

func NewRabbitClient(connectionURL string) (*RabbitClient, error) {
	conn, channel, err := newConnection(connectionURL)
	if err != nil {
		return nil, err
	}

	return &RabbitClient{
		conn:    conn,
		channel: channel,
		url:     connectionURL,
	}, nil
}

func (client *RabbitClient) ConsumeMessages(exchange, routingKey, queueName string) (<-chan amqp.Delivery, error) {
	err := client.channel.ExchangeDeclare(
		exchange,
		"direct",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}
	queue, err := client.channel.QueueDeclare(queueName, true, true, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}
	err = client.channel.QueueBind(queue.Name, routingKey, exchange, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}
	msgs, err := client.channel.Consume(
		queueName,
		"goapp",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume menssages from queue: %v", err)
	}
	return msgs, nil
}

// PublishMessage publishes a message to a specified exchange and binds it to a queue
func (client *RabbitClient) PublishMessage(exchange, routingKey, queueName string, message []byte) error {
	// Ensure the exchange exists before publishing
	err := client.channel.ExchangeDeclare(
		exchange, "direct", true, true, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	// Ensure the queue exists
	_, err = client.channel.QueueDeclare(
		queueName, true, true, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	// Bind the queue to the exchange with the routing key
	err = client.channel.QueueBind(queueName, routingKey, exchange, false, nil)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %v", err)
	}

	// Publish the message to the exchange with the routing key
	err = client.channel.Publish(
		exchange, routingKey, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	return nil
}

func (client *RabbitClient) Close() error {
	err := client.channel.Close()
	if err != nil {
		return fmt.Errorf("failed to close channel: %v", err)
	}
	err = client.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %v", err)
	}
	return nil
}
