package storage

import "blog/types"

type Post struct {
	c *Conn
}

func NewPost(c *Conn) Post {
	return Post{
		c: c,
	}
}

func (p Post) Table(data types.Post) error {
	return p.c.Client.AutoMigrate(&data)

}
