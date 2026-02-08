package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func ConnectToRedis() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	if host == "" {
		host = "redis_queue"
	}
	if port == "" {
		port = "6379"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
}

func GetRedis() *redis.Client {
	return rdb
}