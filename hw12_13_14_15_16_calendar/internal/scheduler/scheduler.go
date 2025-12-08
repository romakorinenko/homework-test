package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/app"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/configs"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/dto"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/rabbitmq"
	"golang.org/x/exp/slog"
)

func MustNewScheduler(
	ctx context.Context,
	scheduler configs.Scheduler,
	logger app.Logger,
	storage app.Storage,
	mq *rabbitmq.RabbitMq,
) gocron.Scheduler {
	sched, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	notificationJob, err := sched.NewJob(
		gocron.CronJob(scheduler.NotificationCron, false),
		gocron.NewTask(func() {
			logger.Info("notification job started", slog.Any("time", time.Now()))
			events, err := storage.GetByPeriod(ctx, time.Now().Add(-16*time.Minute), time.Now())
			if err != nil {
				logger.Error("event receiving err", slog.Any("error", err))
				return
			}

			for _, event := range events {
				notification := &dto.Notification{
					ID:        event.ID,
					Title:     event.Title,
					StartDate: event.StartDate,
					UserID:    event.UserID,
				}
				notificationJSON, err := json.Marshal(notification)
				if err != nil {
					logger.Error("json marshal err", slog.Any("error", err))
					return
				}

				err = mq.Produce(ctx, notificationJSON, "application/json")
				if err != nil {
					logger.Error("produce err", slog.Any("error", err))
					return
				}
				logger.Info("event produced", slog.Any("event", event))
			}
		}),
	)
	if err != nil {
		panic(err)
	}
	logger.Info("notification job created", slog.String("jobID", notificationJob.ID().String()))

	deletionJob, err := sched.NewJob(
		gocron.CronJob(scheduler.EventDeletionCron, false),
		gocron.NewTask(func() {
			logger.Info("deletion job started", slog.Any("time", time.Now()))
			events, err := storage.GetByPeriod(
				ctx,
				time.Now().AddDate(-1, 0, -1),
				time.Now().AddDate(-1, 0, 0),
			)
			if err != nil {
				logger.Error("event receiving err", slog.Any("error", err))
				return
			}

			for _, event := range events {
				err := storage.DeleteByID(ctx, event.ID)
				if err != nil {
					logger.Error("delete err", slog.Any("error", err))
				}
				slog.Info("old event deleted", slog.Any("event", events), slog.Any("event", event))
			}
		}),
	)
	if err != nil {
		panic(err)
	}
	logger.Info("deletion job created", slog.String("jobID", deletionJob.ID().String()))

	gocron.WithStartAt(gocron.WithStartImmediately())
	sched.Start()
	logger.Info("scheduler started")

	return sched
}
