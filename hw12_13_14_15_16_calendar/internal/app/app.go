package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	log     Logger
	storage Storage
}

type Logger interface {
	Info(msg string, data ...interface{})
	Debug(msg string, data ...interface{})
	Warn(msg string, data ...interface{})
	Error(msg string, data ...interface{})
}

type Storage interface {
	io.Closer
	Create(ctx context.Context, event *storage.Event) (*storage.Event, error)
	Update(ctx context.Context, event *storage.Event) error
	DeleteByID(ctx context.Context, id int64) error
	GetByPeriod(ctx context.Context, startData, endData time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{log: logger, storage: storage}
}

func (a *App) CreateEvent(
	ctx context.Context,
	startData, endData time.Time,
	userID int64,
	title string,
) (*storage.Event, error) {
	a.log.Info("CreateEvent - start.",
		slog.Any("startData", startData),
		slog.Any("endData", endData),
		slog.String("title", title),
		slog.Any("userID", userID),
	)

	event, err := a.storage.Create(ctx, &storage.Event{
		StartDate: startData,
		EndDate:   endData,
		Title:     title,
		UserID:    userID,
	})
	if err != nil {
		return nil, fmt.Errorf("create event: %w", err)
	}
	a.log.Info("CreateEvent - end.")

	return event, nil
}

func (a *App) UpdateEvent(
	ctx context.Context,
	id int64,
	startData, endData time.Time,
	userID int64,
	title string,
) error {
	a.log.Info("UpdateEvent - start.",
		slog.Any("id", id),
		slog.Any("startData", startData),
		slog.Any("endData", endData),
		slog.Any("title", title),
	)

	err := a.storage.Update(ctx, &storage.Event{
		ID:        id,
		Title:     title,
		StartDate: startData,
		EndDate:   endData,
		UserID:    userID,
	})
	if err != nil {
		return fmt.Errorf("update event with id '%d': %w", id, err)
	}
	a.log.Info("UpdateEvent - end.")

	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id int64) error {
	a.log.Info("DeleteEvent - start.", slog.Int64("id", id))
	err := a.storage.DeleteByID(ctx, id)
	if err != nil {
		return fmt.Errorf("delete event by id '%d' err: %w", id, err)
	}
	a.log.Info("DeleteEvent - end.")
	return nil
}

func (a *App) GetByPeriod(ctx context.Context, startData, endData time.Time) ([]storage.Event, error) {
	a.log.Info("GetByPeriod - start", slog.Any("startData", startData), slog.Any("endData", endData))

	byPeriod, err := a.storage.GetByPeriod(ctx, startData, endData)
	if err != nil {
		return nil, fmt.Errorf("get event by period err: %w", err)
	}
	a.log.Info("GetByPeriod - end")
	return byPeriod, nil
}

func (a *App) GetLogger() Logger {
	return a.log
}
