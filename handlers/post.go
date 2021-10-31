package handlers

import (
	"blog/storage"
	"blog/types"
	"encoding/json"
	"net/http"
)

type PostHandler struct {
	post storage.Post
}

func NewPostHandler(p storage.Post) PostHandler {
	return PostHandler{
		post: p,
	}
}

func (p PostHandler) Table(w http.ResponseWriter, r *http.Request) {
	var data types.Post
	var err error

	err = p.post.Table(data)

	if err != nil {
		json.NewEncoder(w).Encode("Unable to create table")
	}
	json.NewEncoder(w).Encode(("Table created successfuly"))
}
