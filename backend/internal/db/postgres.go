package db

import (
	"context"
	"fmt"
	"time"

	"github.com/Voltage11/iplatform/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.DatabaseConfig) (*PostgresDB, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &PostgresDB{pool: pool}, nil
}

// Close закрывает пул соединений
func (db *PostgresDB) Close() error {
	db.pool.Close()
	return nil
}

// Pool возвращает пул соединений
func (db *PostgresDB) Pool() *pgxpool.Pool {
	return db.pool
}
