package domain

import (
	"time"

	"github.com/google/uuid"
)

// User пользователь
type User struct {
	ID           uuid.UUID
	Email        string
	FirstName    string
	LastName     string
	PasswordHash string
	IsAdmin      bool
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
