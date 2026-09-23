package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Voltage11/iplatform/internal/appcontext"
	"github.com/Voltage11/iplatform/internal/appmiddleware"
	"github.com/Voltage11/iplatform/internal/handlers/dto"
	"github.com/Voltage11/iplatform/internal/service"
	"github.com/Voltage11/iplatform/internal/utils/httputil"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type authService interface {
	Login(ctx context.Context, email, password, userAgent, clientIP string) (*service.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (*service.TokenPair, error)
	ValidateAccessToken(tokenString string) (uuid.UUID, error)
}

type AuthHandler struct {
	authService    authService
	authMiddleware *appmiddleware.AuthMiddleware
}

func NewAuthHandler(r chi.Router, authService authService, authMiddleware *appmiddleware.AuthMiddleware) {
	authHandler := AuthHandler{
		authService:    authService,
		authMiddleware: authMiddleware,
	}

	// Публичные маршруты
	r.Post("/api/v1/auth/login", authHandler.login)
	r.Post("/api/v1/auth/refresh", authHandler.refresh)
	r.Post("/api/v1/auth/logout", authHandler.logout)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.AuthRequired)
		r.Get("/api/v1/profile", authHandler.profile)
	})

}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		httputil.WriteError(w, err)
		return
	}

	if err := loginRequest.Validate(); err != nil {
		httputil.WriteError(w, err)
		return
	}

	userAgent := r.UserAgent()
	clientIP := r.RemoteAddr

	tokenPair, err := h.authService.Login(r.Context(), loginRequest.Email, loginRequest.Password, userAgent, clientIP)
	if err != nil {
		httputil.WriteError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, dto.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	})
}

func (h *AuthHandler) refresh(w http.ResponseWriter, r *http.Request) {
	var refreshRequest dto.RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&refreshRequest); err != nil {
		httputil.WriteError(w, err)
		return
	}

	if err := refreshRequest.Validate(); err != nil {
		httputil.WriteError(w, err)
		return
	}

	tokenPair, err := h.authService.Refresh(r.Context(), refreshRequest.RefreshToken)
	if err != nil {
		httputil.WriteError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, dto.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	})
}

func (h *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
	var logoutRequest dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&logoutRequest); err != nil {
		httputil.WriteError(w, err)
		return
	}

	if err := logoutRequest.Validate(); err != nil {
		httputil.WriteError(w, err)
		return
	}

	if err := h.authService.Logout(r.Context(), logoutRequest.RefreshToken); err != nil {
		httputil.WriteError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]string{"message": "успешно вышли из системы"})
}

// Profile возвращает информацию о текущем пользователе
func (h *AuthHandler) profile(w http.ResponseWriter, r *http.Request) {
	user := appcontext.GetUserFromContext(r.Context())
	if user == nil {
		httputil.WriteErrorString(w, http.StatusUnauthorized, "пользователь не авторизован")
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]any{
		"user": dto.UserResponse{
			ID:        user.ID.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			IsAdmin:   user.IsAdmin,
		},
	})
}
