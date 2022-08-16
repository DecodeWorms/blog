package pkg

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const cost = 14

func HashPassword(password string) (string, error) {
	p, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", errors.New("unable to hash user password")
	}
	pas := string(p)
	return pas, nil
}

func ComparePassword(hashPass, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(pass))

	if err != nil {
		return errors.New("error validating password")
	}
	return nil
}
