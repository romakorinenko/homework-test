package configs

import (
	"flag"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	configFilePathArgName = "config-file-path"
	defaultConfigFilePath = "configs/config.yaml"
)

func GetAppConfig[C any]() C {
	commandLineFlags := ParseCommandLineFlags()

	var cfg C
	MustLoadConfig(&cfg, commandLineFlags.ConfigFilePath)
	return cfg
}

type CommandLineFlags struct {
	ConfigFilePath string
}

func ParseCommandLineFlags() *CommandLineFlags {
	commandLineFlags := CommandLineFlags{}

	flag.StringVar(&commandLineFlags.ConfigFilePath, configFilePathArgName, defaultConfigFilePath, "a string")

	flag.Parse()

	return &commandLineFlags
}

func MustLoadConfig[C any](config *C, configPath string) {
	slog.Info("loading config file", slog.String("file", configPath))
	mustValidateConfigPath(configPath)

	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		panic(err)
	}
}

func mustValidateConfigPath(configPath string) {
	stat, err := os.Stat(configPath)
	if err != nil {
		panic(err)
	}
	if stat.IsDir() {
		panic("stat cannot be a directory")
	}
}
