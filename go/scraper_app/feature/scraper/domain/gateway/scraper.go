package gateway

import "sales_monitor/scraper_app/feature/scraper/domain/entity"

type Scraper interface {
	GetMarketplaceName() string
	Scrape(url string, cachedProducts *entity.LaterScrapedProducts, wordsToIgnore []string) *entity.ScrapeResult
}

type ScraperFactory interface {
	Get(shopID string) (Scraper, error)
	Close()
}
