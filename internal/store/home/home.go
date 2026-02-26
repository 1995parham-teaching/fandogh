package home

import (
	"context"

	"github.com/1995parham-teaching/fandogh/internal/model"
)

// ListResult contains paginated list of homes with total count.
type ListResult struct {
	Homes []model.Home `json:"homes"`
	Total int64        `json:"total"`
	Skip  int64        `json:"skip"`
	Limit int64        `json:"limit"`
}

// Home stores the home model into the database and S3. we use S3-compatible storage for storing the image files of each home.
type Home interface {
	Set(ctx context.Context, home *model.Home, photos []model.Photo) error
	Get(ctx context.Context, id string) (model.Home, error)
	List(ctx context.Context, skip, limit int64) (ListResult, error)
	Update(ctx context.Context, id string, home model.Home) error
}
