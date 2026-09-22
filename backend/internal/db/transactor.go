package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Transactor управляет транзакциями на уровне сервиса
// Кладёт открытую транзакцию в context, а репозитории через QuerierFrom подхватывают её автоматически
type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// WithinTransaction выполняет fn внутри транзакции
// - fn вернула ошибку  → rollback, ошибка возвращается как есть
// - fn завершилась ок  → commit
func (db *PostgresDB) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// Уже внутри транзакции — просто вызываем fn с тем же ctx.
	if _, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return fn(ctx)
	}

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Rollback при панике — до того, как она улетит выше
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			// Обе ошибки важны: оригинал и причина
			return fmt.Errorf("%w; rollback failed: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
