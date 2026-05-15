package gateway

import (
	"sales_monitor/scraper_app/feature/product/domain/entity"
	"sales_monitor/scraper_app/feature/product/domain/exception"
)

type NotificationPublisher interface {
	SendNotification(notificationTask *entity.NotificationTask) exception.IDomainError
}
