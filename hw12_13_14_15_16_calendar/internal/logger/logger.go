package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(level slog.Level) *Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return &Logger{logger: slog.New(handler)}
}

func (l Logger) Info(msg string, data ...interface{}) {
	l.logger.Info(msg, data...)
}

func (l Logger) Error(msg string, data ...interface{}) {
	l.logger.Error(msg, data...)
}

func (l Logger) Debug(msg string, data ...interface{}) {
	l.logger.Debug(msg, data...)
}

func (l Logger) Warn(msg string, data ...interface{}) {
	l.logger.Warn(msg, data...)
}
