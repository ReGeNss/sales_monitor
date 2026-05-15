package notification

import (
	data "sales_monitor/scraper_app/feature/notification/data/repository"
	"sales_monitor/scraper_app/feature/notification/domain/repository"

	"github.com/redis/go-redis/v9"
)

func NewNotificationFeature(redisClient *redis.Client) repository.NotificationRepository {
	return data.NewNotificationRepository(redisClient)
}