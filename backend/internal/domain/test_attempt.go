package domain

import (
	"time"

	"github.com/google/uuid"
)

type AttemptStatus string

const (
	AttemptStatusInProgress AttemptStatus = "in_progress" // В процессе
	AttemptStatusCompleted  AttemptStatus = "completed"   // Завершена
	AttemptStatusExpired    AttemptStatus = "expired"     // Просрочена/прервана, добавлю позднее период действия прохождения
)

// TestAttempt попытка прохождения темы пользователем
type TestAttempt struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	ThemeID       uuid.UUID
	AttemptNumber int // Порядковый номер попытки
	Status        AttemptStatus
	MaxPoints     int  // Максимальное кол. баллов на момент прохождения, пока фиксирую, т.к. вес вопроса можно менять
	CheckPoints   int  // Аналогично, пока фиксирую из эталона на начала тестирования
	Points        int  // Фактическое значение прохождения
	IsPassed      bool // Зачтено ли
	StartedAt     time.Time
	CompletedAt   *time.Time
}
