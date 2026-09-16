package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
)

// ConfigLog конфигурация лога
type ConfigLog struct {
	Level int
}

// ServerConfig конфигурация сервера
type ServerConfig struct {
	Port            string
	ReadTimeoutSec  time.Duration
	WriteTimeoutSec time.Duration
	IdleTimeoutSec  time.Duration
}

// Config конфигурация приложения
type Config struct {
	Log      ConfigLog
	Server   ServerConfig
	Database DatabaseConfig
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("file .env not load: %v", err)
	}
	// Логгер
	configLog := ConfigLog{
		Level: getEnvInt("LOG_LEVEL", 0),
	}

	// Сервер
	serverConfig := newServerConfig()

	// DB
	databaseConfig, err := newDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки database конфигурации: %v", err)
	}

	return &Config{
		Log:      configLog,
		Server:   serverConfig,
		Database: *databaseConfig,
	}, nil
}

// newServerConfig создание конфига, чтение из переменных окружения
func newServerConfig() ServerConfig {
	return ServerConfig{
		Port:            getEnv("SERVER_PORT", "8080"),
		ReadTimeoutSec:  time.Duration(getEnvInt("SERVER_READ_TIMEOUT_SEC", 10)),
		WriteTimeoutSec: time.Duration(getEnvInt("SERVER_WRITE_TIMEOUT_SEC", 60)),
		IdleTimeoutSec:  time.Duration(getEnvInt("SERVER_IDLE_TIMEOUT_SEC", 60)),
	}
}
