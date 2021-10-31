package handlers

import (
	"blog/storage"
	"blog/types"
	"encoding/json"
	"net/http"
)

type CommentHandler struct {
	comment storage.Comment
}

func NewCommentHandler(com storage.Comment) CommentHandler {
	return CommentHandler{
		comment: com,
	}
}

func (c CommentHandler) Table(w http.ResponseWriter, r *http.Request) {
	var err error
	var data types.Comment
	err = c.comment.Table(data)
	if err != nil {
		json.NewEncoder(w).Encode("Unable to create a table")
	}
	json.NewEncoder(w).Encode("Table created successfully")
}
