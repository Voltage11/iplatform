package service

import (
	"context"
	"time"

	"github.com/Voltage11/iplatform/internal/domain"
	"github.com/Voltage11/iplatform/internal/types/apperr"
	"github.com/google/uuid"
)

type SessionRepo interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByRefreshTokenHash(ctx context.Context, hash string) (*domain.Session, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) (int64, error)
}

type AuthService struct {
	userService *UserService
	sessionRepo SessionRepo
	jwt         *JWTService
}

func NewAuthService(userService *UserService, sessionRepo SessionRepo, jwt *JWTService) *AuthService {
	return &AuthService{
		userService: userService,
		sessionRepo: sessionRepo,
		jwt:         jwt,
	}
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// Login проверяет email/пароль, создаёт сессию и возвращает токены
func (a *AuthService) Login(ctx context.Context, email, password, userAgent, clientIP string) (*TokenPair, error) {
	// 1. Найти пользователя по email
	user, err := a.userService.GetByEmail(ctx, email)
	if err != nil {
		// 1. Проверим на отсутствие пользователя, намеренно не отдаем отсутсвие пользователя если его нет
		if apperr.IsTypeAppError(err, apperr.ErrNotFound) {
			return nil, apperr.NewUnauthorized("неверный email или пароль", nil)
		}
		// Может внутренняя ошибка, выше залогируем
		return nil, err
	}
	// 2. Проверить пароль
	if !a.userService.VerifyPassword(user.PasswordHash, password) {
		return nil, apperr.NewUnauthorized("неверный email или пароль", nil)
	}
	// 3. Проверить активен ли пользователь
	if !user.IsActive {
		return nil, apperr.NewForbidden("учётная запись деактивирована", nil)
	}

	// 4. Сгенерировать токены
	accessToken, err := a.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, apperr.NewInternal("ошибка генерации токена", err)
	}
	refreshToken, err := a.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, apperr.NewInternal("ошибка генерации refresh токена", err)
	}

	// 5. Сохранить сессию
	refreshHash := HashRefreshToken(refreshToken)

	session := &domain.Session{
		ID:               uuid.New(),
		UserID:           user.ID,
		RefreshTokenHash: refreshHash,
		ExpiresAt:        time.Now().UTC().Add(a.jwt.GetRefreshTTL()),
		CreatedAt:        time.Now().UTC(),
		UserAgent:        userAgent,
		ClientIP:         clientIP,
	}
	if err := a.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(a.jwt.GetAccessTTL().Seconds()),
	}, nil
}

// Logout отзывает сессию по refresh токену
func (a *AuthService) Logout(ctx context.Context, refreshToken string) error {
	hashToken := HashRefreshToken(refreshToken)

	session, err := a.sessionRepo.GetByRefreshTokenHash(ctx, hashToken)
	if err != nil {
		return apperr.NewNotFound("сессия не найдена", nil)
	}
	return a.sessionRepo.Revoke(ctx, session.ID)
}

// Refresh выдаёт новую пару токенов по refresh токену
func (a *AuthService) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	hashToken := HashRefreshToken(refreshToken)

	session, err := a.sessionRepo.GetByRefreshTokenHash(ctx, hashToken)
	if err != nil {
		return nil, apperr.NewUnauthorized("недействительный refresh токен", nil)
	}
	if session.RevokedAt != nil {
		return nil, apperr.NewUnauthorized("сессия отозвана", nil)
	}
	if time.Now().UTC().After(session.ExpiresAt) {
		return nil, apperr.NewUnauthorized("refresh токен истёк", nil)
	}

	// Генерируем новый refresh
	newRefreshToken, err := a.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, apperr.NewInternal("ошибка генерации refresh токена", err)
	}
	newRefreshHash := HashRefreshToken(newRefreshToken)

	newSession := &domain.Session{
		ID:               uuid.New(),
		UserID:           session.UserID,
		RefreshTokenHash: newRefreshHash,
		ExpiresAt:        time.Now().UTC().Add(a.jwt.GetRefreshTTL()),
		CreatedAt:        time.Now().UTC(),
		UserAgent:        session.UserAgent,
		ClientIP:         session.ClientIP,
	}
	if err := a.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, err
	}

	// Отзываем старую сессию
	_ = a.sessionRepo.Revoke(ctx, session.ID)

	// Новый access
	accessToken, err := a.jwt.GenerateAccessToken(session.UserID)
	if err != nil {
		return nil, apperr.NewInternal("ошибка генерации токена", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(a.jwt.GetAccessTTL().Seconds()),
	}, nil
}

// ValidateAccessToken валидирует access токен и возвращает userID как строку (для middleware)
func (a *AuthService) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	return a.jwt.ValidateAccessToken(tokenString)
}
