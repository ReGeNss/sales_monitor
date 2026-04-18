package scrapers

import (
	"sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/service/dto"

	"github.com/playwright-community/playwright-go"
)

type Scraper interface {
	GetMarketplaceName() string
	Scrape(
		browser playwright.Browser,
		url string,
		cachedProducts *entity.LaterScrapedProducts,
	) *dto.ScrapeResult
}
