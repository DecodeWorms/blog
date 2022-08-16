package storage

import (
	models "blog/models"
	"context"
)

type Post struct {
	c *Conn
}

func NewPost(context context.Context, c *Conn) Post {
	return Post{
		c: c,
	}
}

func (p Post) Create(ctx context.Context, data models.Post) error {
	return nil
}

func (p Post) PostById(ctx context.Context, id string) (models.Post, error) {
	return models.Post{}, nil
}

func (p Post) UpdatePostById(ctx context.Context, id string) error {
	return nil
}

func (p Post) DeletePostById(ctx context.Context, id string) error {
	return nil
}
