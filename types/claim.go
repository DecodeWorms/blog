package types

import "github.com/dgrijalva/jwt-go"

type Claim struct {
	Username string `json:"username"`
	UserId   uint   `json:"id"`
	jwt.StandardClaims
}
