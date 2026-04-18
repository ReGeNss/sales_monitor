package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sales_monitor/scraper_app/core/env"
	"sales_monitor/scraper_app/shared/product/data/mapper"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/repository"

	"github.com/redis/go-redis/v9"
)

type notificationPublisherImpl struct {
	redisClient *redis.Client
}

func NewNotificationPublisher(redisClient *redis.Client) repository.NotificationPublisher {
	return &notificationPublisherImpl{redisClient: redisClient}
}

func (p *notificationPublisherImpl) SendNotification(notificationTask *entity.NotificationTask) error {
	if notificationTask == nil {
		return fmt.Errorf("notification task is nil")
	}
	if p.redisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	payload, err := json.Marshal(mapper.NotificationTaskToModel(notificationTask))
	if err != nil {
		return err
	}
	return p.redisClient.LPush(context.Background(), env.GetNotificationQueueKey(), string(payload)).Err()
}
