package core

import (
    "context"
    "fmt"

    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var client *redis.Client

func InitRedis(addr string, password string, db int) {
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password, // no password set
        DB:       db,  // use default DB
    })

	pong, err := rdb.Ping(ctx).Result()
	fmt.Println(pong, err)
	if err != nil {
		panic(err)
	}
	client = rdb
}

func GetRedisClient() *redis.Client {
	return client
}