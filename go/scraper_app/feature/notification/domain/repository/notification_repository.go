package repository

import (
	"sales_monitor/scraper_app/feature/notification/domain/entity"
	"sales_monitor/scraper_app/feature/notification/domain/exception"
)

type NotificationRepository interface {
	SendNotification(notificationTask *entity.NotificationTask) exception.IDomainError
}
