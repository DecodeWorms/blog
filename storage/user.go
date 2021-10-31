package storage

import (
	"blog/types"
)

type User struct {
	c *Conn
}

func NewUser(c *Conn) User {
	return User{
		c: c,
	}

}

func (use User) Table(data types.User) error {
	return use.c.Client.AutoMigrate(&data)

}

func (use User) Create(data types.User) error {

	return use.c.Client.Create(&data).Error
}

func (use User) Login(data types.User) error {

	return use.c.Client.First(&data, "username = ? AND password = ?", data.Username, data.Password).Error
}

func (use User) Details(usernme string) (types.User, error) {
	var u types.User
	return u, use.c.Client.First(&u).Error
}
