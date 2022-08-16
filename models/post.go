package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserId   uint64 `gorm:"userid" json:"userid"`
	Title    string `gorm:"title" json:"title"`
	Body     string `gorm:"body"  json:"body"`
	Image    string `gorm:"image" json:"image"`
	Comments []Comment
}
