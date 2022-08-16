package handlers

import (
	"blog/storage"
)

type PostHandler struct {
	post storage.Post
}

func NewPostHandler(p storage.Post) PostHandler {
	return PostHandler{
		post: p,
	}
}
