//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/app"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/configs"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/scheduler"
	internalgrpc "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/server/grpc/pb"
	internalhttp "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/server/http"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	appStorage         *sqlstorage.Storage
	eventServiceClient pb.EventServiceClient
	client             http.Client
)

func TestMain(m *testing.M) {
	slog.Info("start test")
	ctx := context.Background()

	appConfig := configs.GetAppConfig[configs.CalendarConfig]()
	slog.Info("", slog.Any("config", appConfig))

	appLogger := logger.NewLogger(appConfig.Logger.Level)
	appStorage = sqlstorage.NewStorage(ctx, appConfig.Storage.DBString)
	rabbitMq := rabbitmq.NewRabbitMq(appConfig.RabbitMQ)
	err := rabbitMq.Start()
	if err != nil {
		panic(err)
	}
	sched := scheduler.MustNewScheduler(ctx, appConfig.Scheduler, appLogger, appStorage, rabbitMq)

	calendar := app.New(appLogger, appStorage)

	serverHTTP := internalhttp.NewServer(appConfig.HTTP, calendar)

	serverGRPC := internalgrpc.NewServer(appConfig.GRPC.Host, appConfig.GRPC.Port, appLogger, calendar)
	client = http.Client{}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	eventServiceClient = pb.NewEventServiceClient(conn)

	appLogger.Info("calendar is running...")

	go func() {
		_ = serverHTTP.Start()
	}()

	go func() {
		_ = serverGRPC.Start()
	}()

	code := m.Run()
	slog.Info("stop test")

	_ = serverHTTP.Stop(ctx)
	_ = serverGRPC.Stop()
	_ = sched.Shutdown()
	slog.Info("exit with code", slog.Int("code", code))
	os.Exit(code)
}

func TestCalendar_CreateEvent_SuccessCases(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	slog.Info("TestCalendar create event. success cases")
	event := storage.Event{
		Title:     "test event",
		StartDate: now,
		EndDate:   now.Add(20 * time.Second),
		UserID:    1,
	}
	eventBytes, err := json.Marshal(event)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"http://localhost:8080/event/create",
		bytes.NewBuffer(eventBytes),
	)
	require.NoError(t, err)

	resp, err := client.Do(req)
	_ = resp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	events, err := appStorage.GetByPeriod(ctx, now, now.Add(20*time.Second))
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, events[0].Title, "test event")
	require.NotEmpty(t, events[0].ID)

	create, err := eventServiceClient.Create(ctx, &pb.Event{
		Title:     "test event2",
		StartDate: timestamppb.New(now),
		EndDate:   timestamppb.New(now.Add(20 * time.Second)),
		UserId:    1,
	})
	require.NoError(t, err)
	require.Equal(t, create.Title, "test event2")

	events, err = appStorage.GetByPeriod(ctx, now, now.Add(20*time.Second))
	require.NoError(t, err)
	require.Len(t, events, 2)

	for _, event := range events {
		err := appStorage.DeleteByID(ctx, event.ID)
		require.NoError(t, err)
	}
}

func TestCalendar_CreateEvent_FailedCases(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	eventWithoutTitle := storage.Event{
		StartDate: now,
		EndDate:   now.Add(20 * time.Second),
		UserID:    1,
	}
	eventBytes, err := json.Marshal(eventWithoutTitle)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"http://localhost:8080/event/create",
		bytes.NewBuffer(eventBytes),
	)
	require.NoError(t, err)

	resp, err := client.Do(req)
	_ = resp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	events, err := appStorage.GetByPeriod(ctx, now, now.Add(20*time.Second))
	require.NoError(t, err)
	require.Len(t, events, 0)

	create, err := eventServiceClient.Create(ctx, &pb.Event{
		StartDate: timestamppb.New(now),
		EndDate:   timestamppb.New(now.Add(20 * time.Second)),
		UserId:    1,
	})
	require.Error(t, err)
	require.Nil(t, create)

	events, err = appStorage.GetByPeriod(ctx, now, now.Add(20*time.Second))
	require.NoError(t, err)
	require.Len(t, events, 0)
}

func TestCalendar_GetEventsByPeriod(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	slog.Info("TestCalendar get events by periods")
	dayEvent := &storage.Event{
		Title:     "test event",
		StartDate: now,
		EndDate:   now.Add(20 * time.Second),
		UserID:    1,
	}
	dEvent, err := appStorage.Create(ctx, dayEvent)
	require.NoError(t, err)
	require.NotEmpty(t, dEvent.ID)

	weekEvent := &storage.Event{
		Title:     "test event",
		StartDate: now.AddDate(0, 0, 6),
		EndDate:   now.AddDate(0, 0, 6).Add(20 * time.Second),
		UserID:    1,
	}
	wEvent, err := appStorage.Create(ctx, weekEvent)
	require.NoError(t, err)
	require.NotEmpty(t, wEvent.ID)

	monthEvent := &storage.Event{
		Title:     "test event",
		StartDate: now.AddDate(0, 0, 14),
		EndDate:   now.AddDate(0, 0, 14).Add(20 * time.Second),
		UserID:    1,
	}
	mEvent, err := appStorage.Create(ctx, monthEvent)
	require.NoError(t, err)
	require.NotEmpty(t, mEvent.ID)

	reqEvent := &storage.Event{
		StartDate: now.Add(-10 * time.Hour),
	}
	reqEventBytes, err := json.Marshal(reqEvent)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://localhost:8080/event/list/day",
		bytes.NewBuffer(reqEventBytes),
	)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)

	all, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	slog.Info("TestCalendar get events by period", slog.Any("body", string(all)))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, string(all))

	var respEvent []storage.Event
	err = json.Unmarshal(all, &respEvent)
	require.NoError(t, err)
	require.Len(t, respEvent, 1)

	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://localhost:8080/event/list/week",
		bytes.NewBuffer(reqEventBytes),
	)
	require.NoError(t, err)
	respon, err := client.Do(req)
	require.NoError(t, err)

	all, err = io.ReadAll(respon.Body)
	require.NoError(t, err)
	_ = respon.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, string(all))

	err = json.Unmarshal(all, &respEvent)
	require.NoError(t, err)
	require.Len(t, respEvent, 2)

	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://localhost:8080/event/list/month",
		bytes.NewBuffer(reqEventBytes),
	)
	require.NoError(t, err)
	response, err := client.Do(req)
	require.NoError(t, err)

	all, err = io.ReadAll(response.Body)
	require.NoError(t, err)
	_ = response.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, string(all))

	err = json.Unmarshal(all, &respEvent)
	require.NoError(t, err)
	require.Len(t, respEvent, 3)

	pbEvent := &pb.Event{
		StartDate: timestamppb.New(now.Add(-10 * time.Hour)),
	}
	events, err := eventServiceClient.ListByDay(ctx, pbEvent)
	require.NoError(t, err)
	require.Len(t, events.GetEvents(), 1)

	events, err = eventServiceClient.ListByWeek(ctx, pbEvent)
	require.NoError(t, err)
	require.Len(t, events.GetEvents(), 2)

	events, err = eventServiceClient.ListByMonth(ctx, pbEvent)
	require.NoError(t, err)
	require.Len(t, events.GetEvents(), 3)

	for _, event := range events.GetEvents() {
		err := appStorage.DeleteByID(ctx, event.Id)
		require.NoError(t, err)
	}
}
