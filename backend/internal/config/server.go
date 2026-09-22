package config

import (
	"strings"
	"time"
)

type ServerConfig struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	AllowedOrigins []string
}

func newServerConfig() (ServerConfig, error) {
	read, err := getEnvInt("SERVER_READ_TIMEOUT_SEC", 10)
	if err != nil {
		return ServerConfig{}, err
	}
	write, err := getEnvInt("SERVER_WRITE_TIMEOUT_SEC", 60)
	if err != nil {
		return ServerConfig{}, err
	}
	idle, err := getEnvInt("SERVER_IDLE_TIMEOUT_SEC", 60)
	if err != nil {
		return ServerConfig{}, err
	}

	alloweOriginsStr := getEnv("ALLOWED_ORIGONS", "http://localhost:5173,http://127.0.0.1:5173")
	alloweOrigins := strings.Split(alloweOriginsStr, ",")

	return ServerConfig{
		Port:           getEnv("SERVER_PORT", "8080"),
		ReadTimeout:    time.Duration(read) * time.Second,
		WriteTimeout:   time.Duration(write) * time.Second,
		IdleTimeout:    time.Duration(idle) * time.Second,
		AllowedOrigins: alloweOrigins,
	}, nil
}
