package main

import (
	"encoding/json"
	"log"
	"sales_monitor/internal/models"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func processQueue(rdb *redis.Client, queueKey string, db *gorm.DB, messagingClient *messaging.Client) {
	for {
		result, err := rdb.BLPop(ctx, 0*time.Second, queueKey).Result()
		if err != nil {
			log.Printf("redis BLPOP error: %v", err) 
			time.Sleep(2 * time.Second)
			continue
		}

		if len(result) >= 2 {
			payload := result[1]
			var notificationTask models.NotificationTask
			err := json.Unmarshal([]byte(payload), &notificationTask)
			if err != nil {
				log.Printf("error unmarshalling notification task: %v", err)
				continue
			}
			log.Printf("received notification task: %+v", notificationTask)
			sendNotification(db, messagingClient, notificationTask)
		}
	}
}