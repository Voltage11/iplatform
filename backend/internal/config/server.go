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
	RequestTimeout time.Duration
	AllowedOrigins []string
}

func newServerConfig() (ServerConfig, error) {
	read, err := getEnvInt("SERVER_READ_TIMEOUT_SEC", 10)
	if err != nil {
		return ServerConfig{}, err
	}
	write, err := getEnvInt("SERVER_WRITE_TIMEOUT_SEC", 10)
	if err != nil {
		return ServerConfig{}, err
	}
	idle, err := getEnvInt("SERVER_IDLE_TIMEOUT_SEC", 60)
	if err != nil {
		return ServerConfig{}, err
	}
	reqTimeout, err := getEnvInt("SERVER_REQUEST_TIMEOUT_SEC", 15)
	if err != nil {
		return ServerConfig{}, err
	}

	rawOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://127.0.0.1:5173")
	parts := strings.Split(rawOrigins, ",")
	origins := make([]string, 0, len(parts))
	for _, o := range parts {
		if o = strings.TrimSpace(o); o != "" {
			origins = append(origins, o)
		}
	}

	return ServerConfig{
		Port:           getEnv("SERVER_PORT", "8080"),
		ReadTimeout:    time.Duration(read) * time.Second,
		WriteTimeout:   time.Duration(write) * time.Second,
		IdleTimeout:    time.Duration(idle) * time.Second,
		RequestTimeout: time.Duration(reqTimeout) * time.Second,
		AllowedOrigins: origins,
	}, nil
}
