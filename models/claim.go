package models

import "github.com/dgrijalva/jwt-go"

type Claim struct {
	Username string `json:"user_name"`
	UserId   uint   `json:"user_id"`
	Email    string `json:"email"`
	jwt.StandardClaims
}
