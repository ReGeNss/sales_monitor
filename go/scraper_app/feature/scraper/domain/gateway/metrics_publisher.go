package gateway

import "sales_monitor/scraper_app/feature/scraper/domain/entity"

type ScrapingMetrics struct {
	Found   int
	Scraped int
	New     int
	OnSale  int
}

type MetricsPublisher interface {
	Publish(metrics ScrapingMetrics, results map[string]*entity.ScrapingResult)
}
