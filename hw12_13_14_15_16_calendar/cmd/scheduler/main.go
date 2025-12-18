package main

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

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
	appLogger.Info("config received", slog.Any("config", appConfig))

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
	err := rabbitMq.Start()
	if err != nil {
		panic(err)
	}
	newScheduler := scheduler.MustNewScheduler(appCtx, appConfig.Scheduler, appLogger, storage, rabbitMq)

	http.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:              ":8080",
		Handler:           nil,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		_ = server.ListenAndServe()
	}()

	ctx, cancel := signal.NotifyContext(appCtx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-ctx.Done()
	_ = server.Shutdown(appCtx)
	cancel()
	_ = newScheduler.Shutdown()
}
