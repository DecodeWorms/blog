package types

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"username" json:"username,omitempty"`
	Password string `gorm:"password" json:"password,omitempty"`
	Gender   string `gorm:"gender" json:"gender"`
	Location string `gorm:"location"  json:"location,omitempty"`
}
