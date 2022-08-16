package handlers

import (
	"blog/storage"
)

type CommentHandler struct {
	comment storage.Comment
}

func NewCommentHandler(com storage.Comment) CommentHandler {
	return CommentHandler{
		comment: com,
	}
}
