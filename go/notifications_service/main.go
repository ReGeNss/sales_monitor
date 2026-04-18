package main

import (
	"context"
	"log"
	"os"
	"sales_monitor/internal/db"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var ctx = context.Background()

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: .env file not loaded: %v\n", err)
	}

	db.ConnectToRedis()
	db.ConnectToDB()

	messagingClient := initMessaging()

	queueKey := os.Getenv("NOTIFICATION_QUEUE_KEY")
	if queueKey == "" {
		log.Fatalf("missing required env var: NOTIFICATION_QUEUE_KEY")
	}

	processQueue(db.GetRedis(), queueKey, db.GetDB(), messagingClient)
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
