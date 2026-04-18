package gateway

import "sales_monitor/scraper_app/shared/product/domain/entity"

type NotificationPublisher interface {
	SendNotification(notificationTask *entity.NotificationTask) error
}
