package gateway

import "sales_monitor/scraper_app/feature/scraper/domain/entity"

type ResultStorage interface {
	Save(results map[string]*entity.ScrapingResult, categories []string)
}
