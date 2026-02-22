package entity

import (
	"sales_monitor/scraper_app/shared/product/domain/entity"

	"github.com/playwright-community/playwright-go"
)

type ScrapeResult struct {
	Products []*entity.ScrapedProduct
	FoundCount int
	NewCount int
}

type Scraper interface {
	GetMarketplaceName() string
	Scrape(
		browser playwright.Browser,
		url string,
		wordsToIgnore []string,
		cachedProducts *LaterScrapedProducts,
	) *ScrapeResult
}

