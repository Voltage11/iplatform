package appcontext

import (
	"context"

	"github.com/Voltage11/iplatform/internal/domain"
)

type ctxKey int

const userKey ctxKey = iota

func WithUser(ctx context.Context, u *domain.User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func GetUserFromContext(ctx context.Context) *domain.User {
	user, ok := ctx.Value(userKey).(*domain.User)
	if !ok {
		return nil
	}
	return user
}
