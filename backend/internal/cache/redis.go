package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/Voltage11/iplatform/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func New(cfg *config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address(),
		Password: cfg.Password,
		DB:       0, // Используем дефолтную базу
	})

	// Проверяем соединение
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	return &RedisCache{client: client}, nil
}

func (r *RedisCache) Client() *redis.Client {
	return r.client
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
