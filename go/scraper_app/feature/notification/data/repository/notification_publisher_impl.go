package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sales_monitor/scraper_app/core/env"
	"sales_monitor/scraper_app/feature/notification/data/mapper"
	"sales_monitor/scraper_app/feature/notification/domain/entity"
	"sales_monitor/scraper_app/feature/notification/domain/exception"
	"sales_monitor/scraper_app/feature/notification/domain/repository"
	"github.com/redis/go-redis/v9"
)

type notificationRepositoryImpl struct {
	redisClient *redis.Client
}

func NewNotificationRepository(redisClient *redis.Client) repository.NotificationRepository {
	return &notificationRepositoryImpl{redisClient: redisClient}
}

func (p *notificationRepositoryImpl) SendNotification(notificationTask *entity.NotificationTask) exception.IDomainError {
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
