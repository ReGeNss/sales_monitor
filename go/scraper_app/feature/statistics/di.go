package statistics

import (
	"sales_monitor/scraper_app/feature/statistics/data/repository"
	domainRepository "sales_monitor/scraper_app/feature/statistics/domain/repository"
)

func NewStatisticsRepository() domainRepository.StatisticsRepository {
	return repository.NewPrometheusPublisher()
}
