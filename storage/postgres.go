package storage

import (
	"blog/config"
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Conn struct {
	Client *gorm.DB
}

func NewConn(c config.Config, db *gorm.DB) *Conn {
	log.Println("Connecting to DB....")
	var err error

	uri := fmt.Sprintf("host=%s dbname=%s user=%s port=%s", c.DatabaseHost, c.DatabaseName, c.DatabaseUsername, c.DatabasePort)

	db, err = gorm.Open(postgres.Open(uri), &gorm.Config{})

	if err != nil {
		e := errors.New(fmt.Sprintf("Unable to connect to DB %V", err))
		fmt.Println(e)
	}

	log.Println("Connected to DB..")

	return &Conn{
		Client: db,
	}
}
