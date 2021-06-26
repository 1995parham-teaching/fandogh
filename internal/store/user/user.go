package user

import (
	"context"

	"github.com/1995parham/fandogh/internal/model"
)

// User stores and retrieves users.
type User interface {
	Set(ctx context.Context, user model.User) error
	Get(ctx context.Context, email string) (model.User, error)
}
