package apperr

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// NewPostgresError переводит ошибку драйвера pgx в AppError
func NewPostgresError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return NewTimeout("превышено время ожидания запроса к БД", err)
	}
	if errors.Is(err, context.Canceled) {
		return NewTimeout("запрос к БД отменён", err)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return NewNotFound("запись не найдена", err)
	}

	if errors.Is(err, pgx.ErrTooManyRows) {
		return NewInternal("запрос вернул больше одной строки", err)
	}

	// Ошибки PostgreSQL
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return translatePgError(pgErr)
	}

	// Всё остальное — внутренняя ошибка
	return NewInternal("внутренняя ошибка базы данных", err)
}

func translatePgError(pgErr *pgconn.PgError) error {
	switch pgErr.Code {
	// 23505: unique_violation — дубликат уникального ключа
	case pgerrcode.UniqueViolation:
		return NewAlreadyExists("такая запись уже существует", pgErr)

	// 23503: foreign_key_violation — ссылка на несуществующую сущность
	case pgerrcode.ForeignKeyViolation:
		return NewBadRequest("нарушение связности данных", pgErr)

	// 23514: check_violation — данные не прошли CHECK-ограничение
	case pgerrcode.CheckViolation:
		return NewBadRequest("данные не соответствуют ограничениям", pgErr)

	// 23502: not_null_violation — код не заполнил обязательное поле
	case pgerrcode.NotNullViolation:
		return NewInternal("обязательное поле не заполнено", pgErr)

	case pgerrcode.SerializationFailure, pgerrcode.DeadlockDetected:
		return NewInternal("конфликт транзакций, повторите запрос", pgErr)
	// 57014: query_canceled — Postgres сам прервал запрос
	case pgerrcode.QueryCanceled:
		return NewTimeout("запрос к БД был прерван", pgErr)
	}

	// Любой необработанный — внутренняя ошибка
	return NewInternal("внутренняя ошибка базы данных", pgErr)
}
