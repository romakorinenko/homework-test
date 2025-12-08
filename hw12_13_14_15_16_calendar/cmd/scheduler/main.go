package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/app"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/configs"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/scheduler"
	storage2 "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage/sql"
)

func main() {
	appCtx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	appConfig := configs.GetAppConfig[configs.CalendarConfig]()

	appLogger := logger.NewLogger(appConfig.Logger.Level)

	var storage app.Storage
	switch appConfig.Storage.Type {
	case "sql":
		storage = sqlstorage.NewStorage(appCtx, appConfig.Storage.DBString)
		appLogger.Info("using sql storage")
	case "inmemory":
		memorystorage.NewStorage(make(map[int64]*storage2.Event))
		appLogger.Info("using inmemory storage")
	default:
		panic("invalid storage type")
	}

	rabbitMq := rabbitmq.NewRabbitMq(appConfig.RabbitMQ)
	newScheduler := scheduler.MustNewScheduler(appCtx, appConfig.Scheduler, appLogger, storage, rabbitMq)

	ctx, cancel := signal.NotifyContext(appCtx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-ctx.Done()
	cancel()
	_ = newScheduler.Shutdown()
}
