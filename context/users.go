package context

import (
	"context"

	"github.com/jerilcj30/jocp/models"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *models.UserTable) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.UserTable {
	val := ctx.Value(userKey)
	user, ok := val.(*models.UserTable)
	if !ok {
		return nil
	}
	return user
}
