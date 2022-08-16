package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username      string `gorm:"username" json:"user_name,omitempty" validate:"required,min=4"`
	Email         string `gorm:"email" validate:"email"`
	Password      string `gorm:"password" json:"password,omitempty" validate:"required,min=10"`
	Gender        string `gorm:"gender" json:"gender" validate:"required"`
	MaritalStatus string `gorm:"marital_status" validate:"required"`
	Age           int    `gorm:"age" validate:"required"`
	PhoneNumber   string `gorm:"phone_number" validate:"required"`
	Title         string `json:"title" gorm:"title"`
}

type UserInput struct {
	Username    string `gorm:"username" json:"user_name,omitempty" validate:"required,min=4"`
	Email       string `json:"email" gorm:"email" validate:"email"`
	Password    string `json:"password" json:"password,omitempty" validate:"required,min=10"`
	Gender      string `json:"gender" json:"gender" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type Address struct {
	gorm.Model
	UserId     uint   `json:"user_id"`
	State      string `json:"state" validate:"required"`
	LGA        string `json:"lga" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`
}

type AddressInput struct {
	UserId     uint   `gorm:"user_id"`
	State      string `json:"state" validate:"required"`
	LGA        string `json:"lga" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`
}
type Login struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type LoginResponse struct {
	Id            uint   `gorm:"id"`
	Email         string `gorm:"email"`
	Gender        string `gorm:"gender"`
	PhoneNumber   string `gorm"phone_number"`
	MaritalStatus string `gorm:: marital_status"`
	TokenString   string `json:"token_string"`
	Title         string `gorm:"title"`
}
type UserPersonalnfo struct {
	MaritalStatus string `json:"maritai_status"`
	Age           int    `json:"age"`
	Title         string `json:"title" validate:"required"`
}
type StoreOtherUserDataInput struct {
	PersonalInfo UserPersonalnfo `json:"personal_info" validate:"required"`
	Address      AddressInput    `json:"address_input" validate:"required"`
}

type MoreInfo struct {
	UserId        uint
	MaritalStatus string `json:"maritai_status"`
	Age           int    `json:"age"`
	Title         string `json:"title" validate:"required"`
}

type UserDetails struct {
	Id       uint    `json:"id"`
	Username string  `json:"user_name" validate:"required"`
	Email    string  `json:"email" validate:"email"`
	Gender   string  `json:"gender" validate:"required"`
	Address  Address `json:"address"`
}
