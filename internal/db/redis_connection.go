package db

import "github.com/redis/go-redis/v9"

var rdb *redis.Client

func ConnectToRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func GetRedis() *redis.Client {
	return rdb
}