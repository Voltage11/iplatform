package dto

import (
	"strings"

	"github.com/Voltage11/iplatform/internal/types/apperr"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginRequest) Validate() error {
	l.Email = strings.TrimSpace(l.Email)

	if len(l.Email) <= 3 {
		return apperr.NewBadRequest("Невалидный формат email", nil)
	}
	if len(l.Password) == 0 {
		return apperr.NewBadRequest("Пароль не может быть пустым", nil) // было "ну"
	}
	return nil
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r *RefreshRequest) Validate() error {
	r.RefreshToken = strings.TrimSpace(r.RefreshToken)

	if len(r.RefreshToken) == 0 {
		return apperr.NewBadRequest("Токен не может быть пустым", nil)
	}

	return nil
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r *LogoutRequest) Validate() error {
	r.RefreshToken = strings.TrimSpace(r.RefreshToken)

	if len(r.RefreshToken) == 0 {
		return apperr.NewBadRequest("Токен не может быть пустым", nil)
	}

	return nil
}

type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"is_admin"`
}
