package util

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func ExtractEmailAddressFromToken(c *gin.Context, authToken string) (string, error) {
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if err != nil {
		//log.Error("could not parse provided token")
		return "", err
	}

	if !token.Valid {
		//log.Error("Authentication failed " + err.Error())
		return "", errors.Wrap(err, "authentication failed: ")
	}

	claims := jwt.MapClaims{}
	token, err = jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	return fmt.Sprintf("%v", claims["email"]), nil
}

func GetUserIdFromContext(ctx *gin.Context) interface{} {
	v, ok := ctx.Get("user_id")
	if ok {
		return v
	}
	return "no identifier found"
}
