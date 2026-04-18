package env

import (
	"log"
	"os"
)

const (
	jobID = "SCRAPER_JOB_ID"
	notificationQueueKey = "NOTIFICATION_QUEUE_KEY"
)

func getRequiredEnv(key string) (string) {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func GetJobID() string {
	return getRequiredEnv(jobID)
}

func GetNotificationQueueKey() string {
	return getRequiredEnv(notificationQueueKey)
}