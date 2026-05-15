package repository

import "sales_monitor/scraper_app/feature/statistics/domain/entity"

type StatisticsRepository interface {
	Publish(statistics *entity.ScrapingStatistics)
}
