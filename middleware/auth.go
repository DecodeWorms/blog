package middleware

import (
	"blog/auth"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ValidateToken(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")

	if tokenString == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token not available at the http header"})
		ctx.Abort()
		return
	}

	strAr := strings.Split(tokenString, " ")
	if len(strAr) == 2 {
		err := auth.ValidatingToken(strAr[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			fmt.Println(err)
			ctx.Abort()
			return
		}

	}
	if err := EmbedClaimsInContext(ctx, strAr[1]); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

}

func EmbedClaimsInContext(c *gin.Context, authToken string) error {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(
		authToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})

	if err != nil {
		return err
	}

	// save decoded claims in a Context
	c.Set("user_id", claims["user_id"])

	return nil
}
