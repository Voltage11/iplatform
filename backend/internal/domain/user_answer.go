package domain

import "github.com/google/uuid"

type UserAnswer struct {
	ID            uuid.UUID
	TestAttemptID uuid.UUID
	QuestionID    uuid.UUID
	AnswerIDs     []uuid.UUID
	IsCorrect     bool
	Points        int
}
