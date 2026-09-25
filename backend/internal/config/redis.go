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
	return RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
	}
}
