package domain

import (
	"time"

	"github.com/google/uuid"
)

type Theme struct {
	ID           uuid.UUID
	Name         string
	Description  string
	IsActive     bool
	CreatedBy    uuid.UUID
	MaxPoints    int //Макс возможное кол баллоы
	CheckPoints  int //Кол баллов для зачета
	ImagePath    string
	AttemptCount int // Количество попыток пройти тестирование
	CreatedAt    time.Time
}
