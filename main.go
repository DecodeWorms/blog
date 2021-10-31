package main

import (
	"blog/config"
	"blog/handlers"
	"blog/storage"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var c *storage.Conn
var user storage.User
var post storage.Post
var comment storage.Comment

var r *storage.RedisClient

var u handlers.UserHandler

func init() {
	_ = godotenv.Load()
	host := os.Getenv("DATABASE_HOST")
	username := os.Getenv("DATABASE_USERNAME")
	port := os.Getenv("DATABASE_PORT")
	dbName := os.Getenv("DATABASE_NAME")
	//redisPort := os.Getenv("REDIS_PORT")

	var db *gorm.DB

	cfg := config.Config{
		DatabaseHost:     host,
		DatabaseName:     dbName,
		DatabasePort:     port,
		DatabaseUsername: username,
	}
	//conection to DB
	c = storage.NewConn(cfg, db)
	user = storage.NewUser(c)
	post = storage.NewPost(c)
	comment = storage.NewComment(c)
}
func main() {

	u := handlers.NewUserHandler(user)
	p := handlers.NewPostHandler(post)
	c := handlers.NewCommentHandler(comment)
	router := mux.NewRouter()
	router.HandleFunc("/user/table", u.Table).Methods("POST")
	router.HandleFunc("/user/create", u.Create).Methods("POST")
	router.HandleFunc("/user/login", u.Login).Methods("GET")
	router.HandleFunc("/user/details/{username}", u.Details).Methods("GET")

	router.HandleFunc("/post/table", p.Table).Methods("POST")

	router.HandleFunc("/comment/table", c.Table).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))

}
