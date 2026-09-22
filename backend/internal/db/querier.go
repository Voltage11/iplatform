package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Querier — общий интерфейс для pgxpool.Pool и pgx.Tx
// Позволяет репозиториям работать внутри транзакции или без неё
type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// txKey — приватный ключ для хранения транзакции в context
type txKey struct{}

// QuerierFrom возвращает источник запросов:
//   - транзакцию из контекста, если она там есть
//   - сам пул — если транзакции нет
func QuerierFrom(ctx context.Context, pool *pgxpool.Pool) Querier {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return pool
}
