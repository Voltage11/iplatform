package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type LogConfig struct {
	Level int
}

type Config struct {
	Log      LogConfig
	Server   ServerConfig
	Database DatabaseConfig
}

func New() (*Config, error) {
	_ = godotenv.Load()

	logCfg, err := newLogConfig()
	if err != nil {
		return nil, fmt.Errorf("log config: %w", err)
	}
	srv, err := newServerConfig()
	if err != nil {
		return nil, fmt.Errorf("server config: %w", err)
	}
	db, err := newDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("database config: %w", err)
	}

	return &Config{
		Log:      logCfg,
		Server:   srv,
		Database: db,
	}, nil
}

func newLogConfig() (LogConfig, error) {
	level, err := getEnvInt("LOG_LEVEL", 0)
	if err != nil {
		return LogConfig{}, err
	}
	return LogConfig{Level: level}, nil
}
