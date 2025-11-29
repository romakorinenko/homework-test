package rabbitmq

import (
	"context"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/configs"
)

type RabbitMq struct {
	config configs.RabbitMQ
	conn   *amqp.Connection
	ch     *amqp.Channel
	queue  amqp.Queue
}

func NewRabbitMq(config configs.RabbitMQ) *RabbitMq {
	return &RabbitMq{
		config: config,
	}
}

func (r *RabbitMq) Start() error {
	conn, err := amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%d/", r.config.User, r.config.Password, r.config.Host, r.config.Port),
	)
	if err != nil {
		return err
	}
	r.conn = conn

	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	r.ch = ch

	q, err := ch.QueueDeclare(
		r.config.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	r.queue = q

	return nil
}

func (r *RabbitMq) Stop() {
	_ = r.conn.Close()
	_ = r.ch.Close()
}

func (r *RabbitMq) Produce(ctx context.Context, msg []byte, contentType string) error {
	if r.ch == nil {
		return errors.New("channel is nil")
	}

	return r.ch.PublishWithContext(
		ctx,
		"",
		r.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        msg,
		},
	)
}

func (r *RabbitMq) Consume(ctx context.Context) (<-chan amqp.Delivery, error) {
	if r.ch == nil {
		return nil, errors.New("channel is nil")
	}

	messages, err := r.ch.ConsumeWithContext(
		ctx,
		r.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
