package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Username string `gorm:"usersame" json:"username"`
	PostId   uint64 `gorm:"postid" json:"postid"`
	Comment  string `gorm:"comment" json:"comment"`
}
