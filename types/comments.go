package types

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	PostId  uint64 `gorm:"postid" json:"postid"`
	UserId  uint64 `gorm:"userid" json:"userid"`
	Comment string `gorm:"comment" json:"comment"`
}
