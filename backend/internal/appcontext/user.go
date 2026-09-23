package appcontext

import (
	"context"

	"github.com/Voltage11/iplatform/internal/domain"
	"github.com/Voltage11/iplatform/internal/service"
)

func GetUserFromContext(ctx context.Context) *domain.User {
	user, ok := ctx.Value(service.UserContextKey).(*domain.User)
	if !ok {
		return nil
	}
	return user
}
