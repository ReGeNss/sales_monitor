package scrapers

import (
	"sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/service/dto"
)

type Scraper interface {
	GetMarketplaceName() string
	Scrape(
		url string,
		cachedProducts *entity.LaterScrapedProducts,
	) *dto.ScrapeResult
}
