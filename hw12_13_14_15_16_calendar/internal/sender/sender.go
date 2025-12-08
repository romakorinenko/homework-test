package sender

import (
	"context"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/app"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/rabbitmq"
)

type Sender struct {
	queue  *rabbitmq.RabbitMq
	logger app.Logger
}

func NewSender(queue *rabbitmq.RabbitMq, logger app.Logger) Sender {
	return Sender{
		queue:  queue,
		logger: logger,
	}
}

func (s *Sender) Run(ctx context.Context) error {
	ch, err := s.queue.Consume(ctx)
	if err != nil {
		return err
	}
	for {
		select {
		case msg := <-ch:
			s.logger.Info(string(msg.Body))
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *Sender) Start() error {
	return s.queue.Start()
}

func (s *Sender) Stop() {
	s.queue.Stop()
}
