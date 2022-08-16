package storage

import (
	models "blog/models"
	"context"
)

type Comment struct {
	c *Conn
}

func NewComment(context context.Context, c *Conn) Comment {
	return Comment{
		c: c,
	}
}

func (c Comment) Table(data models.Comment) error {
	return c.c.Client.AutoMigrate(&data)
}
