package home

import (
	"context"

	"github.com/1995parham/fandogh/internal/model"
)

// Home stores the home model into the database and minio. we use minio for storing the image files of each home.
type Home interface {
	Set(ctx context.Context, home model.Home) (string, error)
	Get(ctx context.Context, id string) (model.Home, error)
}
