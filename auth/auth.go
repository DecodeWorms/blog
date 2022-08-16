package auth

import (
	"blog/models"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(email, usernme string, id uint) (string, error) {
	//exp time for token
	expTime := time.Now().Add(time.Hour * 1)

	//claim or payload for token
	clm := &models.Claim{
		Username: usernme,
		UserId:   id,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}

	//get instance of *jwt.Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, clm)

	// create token string
	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func ValidatingToken(tokenString string) error {
	// parsing token and claims
	token, err := jwt.ParseWithClaims(tokenString, &models.Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if err != nil {
		fmt.Println(err)
		return errors.New("unable to validate token")
	}
	// extracting claims

	claim, ok := token.Claims.(*models.Claim)
	if !ok {
		return errors.New("unable to extract claims")
	}

	if err := claim.Valid(); err != nil {
		return errors.New("token expired")
	}
	return nil

}
