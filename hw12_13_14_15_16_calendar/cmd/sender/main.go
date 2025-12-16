package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/configs"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/sender"
)

func main() {
	appCtx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	appConfig := configs.GetAppConfig[configs.CalendarConfig]()

	appLogger := logger.NewLogger(appConfig.Logger.Level)
	appLogger.Info("config received", slog.Any("config", appConfig))

	rabbitMq := rabbitmq.NewRabbitMq(appConfig.RabbitMQ)
	newSender := sender.NewSender(rabbitMq, appLogger)
	err := newSender.Start()
	if err != nil {
		panic(err)
	}
	err = newSender.Run(appCtx)
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(appCtx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-ctx.Done()
	cancel()
	newSender.Stop()
}
