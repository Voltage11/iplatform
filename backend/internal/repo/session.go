package repo

import (
	"context"

	"github.com/Voltage11/iplatform/internal/domain"
	"github.com/Voltage11/iplatform/internal/types/apperr"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type sessionRepo struct {
	pool *pgxpool.Pool
}

func NewSessionRepo(pool *pgxpool.Pool) *sessionRepo {
	return &sessionRepo{
		pool: pool,
	}
}

func (s *sessionRepo) Create(ctx context.Context, session *domain.Session) error {
	sql := `
        INSERT INTO sessions (id, user_id, refresh_token_hash, expires_at, created_at, user_agent, client_ip)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	result, err := s.pool.Exec(ctx, sql, session.ID,
		session.UserID,
		session.RefreshTokenHash,
		session.ExpiresAt,
		session.CreatedAt,
		session.UserAgent,
		session.ClientIP)
	if err != nil {
		return apperr.NewPostgresError(err)
	}

	if result.RowsAffected() == 0 {
		return apperr.NewInternal("сессия не создана", apperr.NewNotFound("Сессия для пользователя не создана", nil))
	}

	return nil
}

func (r *sessionRepo) GetByRefreshTokenHash(ctx context.Context, hash string) (*domain.Session, error) {
	query := `
        SELECT id, user_id, refresh_token_hash, expires_at, created_at, revoked_at, user_agent, client_ip
        FROM sessions
        WHERE refresh_token_hash = $1
    `
	var s domain.Session
	err := r.pool.QueryRow(ctx, query, hash).Scan(
		&s.ID, &s.UserID, &s.RefreshTokenHash, &s.ExpiresAt, &s.CreatedAt, &s.RevokedAt,
		&s.UserAgent, &s.ClientIP,
	)
	if err != nil {
		return nil, apperr.NewPostgresError(err)
	}
	return &s, nil
}

func (r *sessionRepo) Revoke(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE sessions SET revoked_at = NOW() WHERE id = $1 AND revoked_at IS NULL`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return apperr.NewPostgresError(err)
	}
	if result.RowsAffected() == 0 {
		return apperr.NewNotFound("сессия не найдена или уже отозвана", nil)
	}
	return nil
}

func (r *sessionRepo) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE sessions SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := r.pool.Exec(ctx, query, userID)
	return apperr.NewPostgresError(err)
}

func (r *sessionRepo) DeleteExpired(ctx context.Context) (int64, error) {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`
	result, err := r.pool.Exec(ctx, query)
	if err != nil {
		return 0, apperr.NewPostgresError(err)
	}

	return result.RowsAffected(), nil
}
