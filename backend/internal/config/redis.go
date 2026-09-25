package config

import (
	"net"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func (c RedisConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func newRedisConfig() RedisConfig {
	// Сначала проверяем локальную переменную, если её нет — берем докер
	password := getEnv("REDIS_PASSWORD", "")
	if password == "" {
		password = getEnv("DOCKER_REDIS_PASSWORD", "")
	}

	return RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: password,
	}
}
