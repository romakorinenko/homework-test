package main

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/configs"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/sender"
)

func main() {
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
	_ = server.Shutdown(appCtx)
	cancel()
	newSender.Stop()
}
