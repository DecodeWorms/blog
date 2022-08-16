package features

import (
	"blog/models"
	"context"
)

type PostServices interface {
	Create(ctx context.Context, data models.Post) error
	// PostById(ctx context.Context, id string) (models.Post, error)
	// UpdatePostById(ctx context.Context, id string, data models.Post) error
	// DeletePostById(ctx context.Context, id string) error
}
