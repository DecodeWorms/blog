package storage

import (
	"blog/models"
	"context"
	"time"
)

type UserServices interface {
	Automigrate() error
	Create(ctx context.Context, input models.UserInput) error
	CreateAddress(ctx context.Context, input models.AddressInput) error
	StoreOtherUserData(ctx context.Context, input models.MoreInfo) error
	Address(ctx context.Context, userId uint) (*models.Address, error)
	UserById(ctx context.Context, id string) (*models.User, error)
	UserByEmail(ctx context.Context, email string) (*models.User, error)
	UserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
	UpdateUserById(ctx context.Context, id string, data models.User) error
	DeleteUserById(ctx context.Context, id string) error
}

type User struct {
	c *Conn
}

func NewUser(ctx context.Context, c *Conn) User {
	return User{
		c: c,
	}

}

func (u User) Automigrate() error {
	var data models.User
	return u.c.Client.AutoMigrate(&data)

}

func (u User) Create(ctx context.Context, input models.UserInput) error {
	var da models.User

	da.CreatedAt, da.UpdatedAt = time.Now(), time.Now()

	d := &models.User{
		Username:    input.Username,
		Email:       input.Email,
		Password:    input.Password,
		Gender:      input.Gender,
		PhoneNumber: input.PhoneNumber,
	}
	return u.c.Client.Create(&d).Error
}

func (u User) StoreOtherUserData(ctx context.Context, input models.MoreInfo) error {
	var d models.User
	d.UpdatedAt = time.Now()

	return u.c.Client.Model(&models.User{}).Where("id = ?", input.UserId).Updates(models.User{
		MaritalStatus: input.MaritalStatus,
		Age:           input.Age,
		Title:         input.Title,
	}).Error

}

func (u User) CreateAddress(ctx context.Context, input models.AddressInput) error {
	data := &models.Address{
		UserId:     input.UserId,
		State:      input.State,
		LGA:        input.LGA,
		PostalCode: input.PostalCode,
	}
	return u.c.Client.Create(data).Error
}

func (u User) Address(ctx context.Context, userId uint) (*models.Address, error) {
	var add *models.Address

	return add, u.c.Client.Where("user_id = ?", userId).First(&add).Error
}

func (u User) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	var m models.User
	if err := u.c.Client.Where("email = ?", email).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (u User) UserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	var us models.User
	return &us, u.c.Client.Where("phone_number = ?", phoneNumber).First(&us).Error

}

func (u User) UserById(ctx context.Context, id string) (*models.User, error) {
	var m models.User
	return &m, u.c.Client.Where("id = ?", id).First(&m).Error

}

func (u User) UpdateUserById(ctx context.Context, id string, data models.User) error {
	return nil
}

func (u User) DeleteUserById(ctx context.Context, id string) error {
	return nil
}
