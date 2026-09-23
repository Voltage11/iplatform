package domain

import (
	"time"

	"github.com/google/uuid"
)

// Session сессии пользователей
type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	RefreshTokenHash string
	ExpiresAt        time.Time  // Время окончания сессии
	RevokedAt        *time.Time // Время отзыва сессии
	CreatedAt        time.Time
	UserAgent        string
	ClientIP         string
}
