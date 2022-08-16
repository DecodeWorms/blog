package server

import (
	"blog/handlers"
	"blog/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserServer struct {
	user handlers.UserHandler
}

func NewUserServer(u handlers.UserHandler) UserServer {
	return UserServer{
		user: u,
	}
}

func (u UserServer) AutoMigrate() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := u.user.AutoMigrate()
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))

		}

		cc := c.Request

		select {
		case <-time.After(3 * time.Second):
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Gateway time out"})
		case <-cc.Context().Done():
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gateway time out"})
		}

	}
}

func (u UserServer) Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data models.UserInput

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		if err := u.user.Signup(ctx, data); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"response ": "User datas saved successfully"})

	}
}

func (u UserServer) StoreOtherUserData() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data models.StoreOtherUserDataInput

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		if err := u.user.StoreOtherUserData(ctx, data); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"response ": "Other User data saved successfully"})

	}

}

func (u UserServer) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var lgn models.Login

		if err := ctx.ShouldBindJSON(&lgn); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		loginResponse, err := u.user.Login(ctx, lgn)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": loginResponse})

	}
}

func (u UserServer) UserDetails() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.Query("email")
		user, err := u.user.UserDetails(ctx, email)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.JSON(200, gin.H{"datas": user})
	}
}

func (u UserServer) UserById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Query("id")
		res, err := u.user.UserById(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"user": res})
	}
}
