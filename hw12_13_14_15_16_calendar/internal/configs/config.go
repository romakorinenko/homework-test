package configs

import (
	"log/slog"
	"time"
)

type CalendarConfig struct {
	Logger  Logger  `yaml:"logger"`
	Storage Storage `yaml:"storage"`
	HTTP    HTTP    `yaml:"http"`
	GRPC    GRPC    `yaml:"grpc"`
}

type Logger struct {
	Level slog.Level `yaml:"level"`
}

type Storage struct {
	Type     string `yaml:"type"`
	DBString string `yaml:"dbString"`
}

type HTTP struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type GRPC struct {
	Host    string        `yaml:"host"`
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}
