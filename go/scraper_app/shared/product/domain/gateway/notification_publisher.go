package gateway

import (
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/exception"
)

type NotificationPublisher interface {
	SendNotification(notificationTask *entity.NotificationTask) exception.IDomainError
}
