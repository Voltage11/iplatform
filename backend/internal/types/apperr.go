package types

import (
	"errors"
	"fmt"
)

// ErrType — кастомная ошибка
type ErrType string

const (
	ErrNotFound      ErrType = "NOT_FOUND"
	ErrBadRequest    ErrType = "BAD_REQUEST"
	ErrAlreadyExists ErrType = "ALREADY_EXISTS"
	ErrUnauthorized  ErrType = "UNAUTHORIZED"
	ErrForbidden     ErrType = "FORBIDDEN"
	ErrInternal      ErrType = "INTERNAL"
)

// AppError — ошибка с категорией
type AppError struct {
	Type    ErrType `json:"type"`
	Message string  `json:"message"`
	Err     error   `json:"-"`
}

// Error реализует интерфейс error
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// Unwrap возвращает оригинальную ошибку
func (e *AppError) Unwrap() error { return e.Err }

// Is позволяет сравнивать ошибки по типу:
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	return ok && t.Type == e.Type
}

// Конструкторы
func NewNotFound(msg string, err error) error {
	return &AppError{Type: ErrNotFound, Message: msg, Err: err}
}

func NewAlreadyExists(msg string, err error) error {
	return &AppError{Type: ErrAlreadyExists, Message: msg, Err: err}
}

func NewBadRequest(msg string, err error) error {
	return &AppError{Type: ErrBadRequest, Message: msg, Err: err}
}

func NewUnauthorized(msg string, err error) error {
	return &AppError{Type: ErrUnauthorized, Message: msg, Err: err}
}

func NewForbidden(msg string, err error) error {
	return &AppError{Type: ErrForbidden, Message: msg, Err: err}
}

func NewInternal(msg string, err error) error {
	return &AppError{Type: ErrInternal, Message: msg, Err: err}
}

// AsAppError извлекает *AppError из цепочки ошибок.
func AsAppError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
