package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/app"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/configs"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/server/http"
	storage2 "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage/sql"
)

func main() {
	printVersion()
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

	calendar := app.New(appLogger, storage)

	serverHTTP := internalhttp.NewServer(appConfig.HTTP, calendar)
	serverGRPC := internalgrpc.NewServer(appConfig.GRPC.Host, appConfig.GRPC.Port, appLogger, calendar)

	appLogger.Info("calendar is running...")

	go func() {
		_ = serverHTTP.Start()
	}()

	go func() {
		_ = serverGRPC.Start()
	}()

	ctx, cancel := signal.NotifyContext(appCtx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-ctx.Done()
	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := serverHTTP.Stop(ctx); err != nil {
		appLogger.Error("failed to stop http serverHTTP: " + err.Error())
	}
	_ = storage.Close()

	if err := serverGRPC.Stop(); err != nil {
		appLogger.Error("failed to stop grpc serverHTTP: " + err.Error())
	}
}
