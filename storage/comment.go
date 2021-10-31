package storage

import "blog/types"

type Comment struct {
	c *Conn
}

func NewComment(c *Conn) Comment {
	return Comment{
		c: c,
	}
}

func (c Comment) Table(data types.Comment) error {
	return c.c.Client.AutoMigrate(&data)
}
