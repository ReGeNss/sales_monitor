package env

import (
	"log"
	"os"
)

const (
	jobID = "SCRAPER_JOB_ID"
)

func getRequiredEnv(key string) (string) {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func getJobID() string {
	return getRequiredEnv(jobID)
}