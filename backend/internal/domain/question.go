package domain

import (
	"time"

	"github.com/google/uuid"
)

type QuestionType string

const (
	QuestionTypeSingle QuestionType = "single" // Один вариант ответа
	QuestionTypeMulti  QuestionType = "multi"  // Более одного варианта ответа
)

type Question struct {
	ID           uuid.UUID
	ThemeID      uuid.UUID
	Name         string
	Points       int // Кол баллов за верный ответ
	SortOrder    int
	QuestionType QuestionType
	CreatedAt    time.Time
}
