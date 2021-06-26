package user

import (
	"context"

	"github.com/1995parham/fandogh/internal/model"
)

// User stores and retrieves users.
type URL interface {
	Set(ctx context.Context, user model.User) string
	Get(ctx context.Context, email string) (model.User, error)
}
