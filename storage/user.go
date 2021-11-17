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

type Result struct {
	Username string
	Gender   string
	Location string
	Title    string
	Body     string
	Image    string
}

func (use User) Automigrate(data types.User) error {
	return use.c.Client.AutoMigrate(&data)
}

func (use User) Create(data types.User) error {

	return use.c.Client.Create(&data).Error
}

func (use User) Login(data types.User) error {

	return use.c.Client.First(&data, "username = ? AND password = ?", data.Username, data.Password).Error
}

func (use User) MyProfiles(username string) ([]types.User, error) {
	var u []types.User
	//return u, use.c.Client.First(&u, "username = ?", username).Error
	//return u, use.c.Client.Model(types.User{}).Select("users.username,users.gender,users.location,posts.title,posts.body,posts.image").Joins("left join posts on posts.user_id = users.id").Where("users.username = ?", username).Scan(&u).Error
	return u, use.c.Client.Preload("Posts").Find(&u, "username = ?", username).Error
}

func (use User) Posts() ([]types.Post, error) {
	var d []types.Post
	return d, use.c.Client.Preload("Comments").Find(&d).Error
}

func (user User) UpdateName(username string, data types.User) error {
	return user.c.Client.Model(&types.User{}).Where("username = ?", username).Update("username", data.Username).Error
}

func (user User) Post(userId uint64, data types.Post) error {
	post := types.Post{
		UserId: uint64(userId),
		Title:  data.Title,
		Body:   data.Body,
		Image:  data.Image,
	}

	return user.c.Client.Create(&post).Error
}

func (user User) Coment(username string, data types.Comment) error {
	cmt := types.Comment{
		Username: username,
		PostId:   uint64(data.PostId),
		Comment:  data.Comment,
	}
	return user.c.Client.Create(&cmt).Error
}
