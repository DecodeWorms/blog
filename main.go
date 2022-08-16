package main

import (
	"blog/config"
	"blog/handlers"
	"blog/middleware"
	"blog/server"
	"blog/storage"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var c *storage.Conn
var user storage.User
var post storage.Post
var comment storage.Comment

var userServer server.UserServer

var userhandler handlers.UserHandler

// var r *storage.RedisClient

// var u handlers.UserHandler

func init() {
	_ = godotenv.Load()
	host := os.Getenv("DATABASE_HOST")
	username := os.Getenv("DATABASE_USERNAME")
	port := os.Getenv("DATABASE_PORT")
	dbName := os.Getenv("DATABASE_NAME")
	var db *gorm.DB

	cfg := config.Config{
		DatabaseHost:     host,
		DatabaseName:     dbName,
		DatabasePort:     port,
		DatabaseUsername: username,
	}
	//conection to DB
	c = storage.NewConn(cfg, db)

	//services layer
	initContext := &gin.Context{}
	user = storage.NewUser(initContext, c)
	post = storage.NewPost(initContext, c)
	comment = storage.NewComment(initContext, c)

	//handler layer
	userhandler = handlers.NewUserHandler(user)

	//server layer
	userServer = server.NewUserServer(userhandler)

}
func main() {
	router := gin.Default()
	router.POST("/user/table", userServer.AutoMigrate())
	router.POST("/user/create", userServer.Signup())
	router.POST("/user/login", userServer.Login())

	//Middleware validates if incoming http request has valid token
	secured := router.Group("/user").Use(middleware.ValidateToken)
	{
		secured.GET("/details", userServer.UserDetails())
		secured.PUT("/kyc", userServer.StoreOtherUserData())
		secured.GET("detail2", userServer.UserById())
	}
	router.Run(":3000")

}
