package types

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"username" json:"username,omitempty" validate:"required,min=11"`
	Password string `gorm:"password" json:"password,omitempty" validate:"required,min=10,passwd"`
	Gender   string `gorm:"gender" json:"gender" validate:"required,min=4"`
	Location string `gorm:"location"  json:"location,omitempty" validate:"required"`
	Posts    []Post
}
