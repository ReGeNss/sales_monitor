package main

import (
	"context"
	"log"
	"os"
	"sales_monitor/internal/db"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/option"
)

var ctx = context.Background()

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v\n", err)
	}

	db := db.GetDB()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()

	messagingClient := initMessaging()

	queueKey := os.Getenv("NOTIFICATIONS_QUEUE_KEY")

	processQueue(rdb, queueKey, db, messagingClient)
}


func initMessaging() *messaging.Client {
	opt := option.WithCredentialsFile("firebase_key.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error initializing firebase messaging client: %v\n", err)
	}

	return messagingClient
}
