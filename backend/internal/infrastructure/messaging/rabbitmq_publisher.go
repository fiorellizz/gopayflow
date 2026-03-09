package messaging

import (
	"context"
	"encoding/json"

	"github.com/streadway/amqp"

	"github.com/fiorellizz/gopayflow/internal/domain"
)

type RabbitMQPublisher struct {
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQPublisher(channel *amqp.Channel, queue string) *RabbitMQPublisher {
	return &RabbitMQPublisher{
		channel: channel,
		queue:   queue,
	}
}

func (p *RabbitMQPublisher) PublishOrderCreated(ctx context.Context, order *domain.Order) error {

	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"",
		p.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}