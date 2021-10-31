package storage

import (
	"fmt"

	"github.com/go-redis/redis/v7"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cl *redis.Client, ad string) *RedisClient {

	cl = redis.NewClient(&redis.Options{
		Addr: ad,
	})
	res, err := cl.Ping().Result()

	if err != nil {
		panic(err)
	}
	fmt.Println(res, "connected")

	c := RedisClient{
		Client: cl,
	}
	return &c

}
