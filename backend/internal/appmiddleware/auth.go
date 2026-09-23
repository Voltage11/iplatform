package appmiddleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Voltage11/iplatform/internal/appcontext"
	"github.com/Voltage11/iplatform/internal/domain"
	"github.com/Voltage11/iplatform/internal/utils/httputil"
	"github.com/google/uuid"
)

type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type JwtService interface {
	ValidateAccessToken(tokenString string) (uuid.UUID, error)
}

// AuthMiddleware мидлеваре для обязательной авторизации и нет
type AuthMiddleware struct {
	userService UserService
	jwtService  JwtService
}

// NewAuthMiddleware Конструктор миддлеваре аутентификации
func NewAuthMiddleware(userService UserService, jwtService JwtService) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
		jwtService:  jwtService,
	}
}

// ExtractUser - Первично провеяет, есть ли токен пользователя в заголовках, если есть и валидный, то
// кладем, если отсутствует или не валидный, то не кладем, но пропускаем дальше
func (a *AuthMiddleware) ExtractUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.SplitN(tokenStr, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			next.ServeHTTP(w, r)
			return
		}
		token := parts[1]

		userID, err := a.jwtService.ValidateAccessToken(token)
		// Не валидный токен, не получили userID
		if err != nil || userID == uuid.Nil {
			next.ServeHTTP(w, r)
			return
		}

		// Если валидный пользователь, проверим его активность
		user, err := a.userService.GetByID(r.Context(), userID)
		if err != nil || user == nil {
			next.ServeHTTP(w, r)
			return
		}

		if !user.IsActive {
			next.ServeHTTP(w, r)
			return
		}

		// Если пользователь верийицирован и активный, то положим в контекст, по идее позже верификацию вынесу в отдельнуй метод, т.к. условия может быть больше
		ctxWithUser := appcontext.WithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctxWithUser))

	})
}

// AuthRequired - обязательная верификация пользователя для закрытых ручек
func (a *AuthMiddleware) AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := appcontext.GetUserFromContext(r.Context())
		if user == nil {
			httputil.WriteErrorString(w, http.StatusUnauthorized, "требуется авторизация")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthAdmin - обязательная верификация пользователя только для админов
func (a *AuthMiddleware) AuthAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := appcontext.GetUserFromContext(r.Context())
		if user == nil {
			httputil.WriteErrorString(w, http.StatusUnauthorized, "требуется авторизация")
			return
		}

		if !user.IsAdmin {
			httputil.WriteErrorString(w, http.StatusForbidden, "требуется авторизация администратора")
			return
		}

		next.ServeHTTP(w, r)
	})
}
