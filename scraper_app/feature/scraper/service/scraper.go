package scraper

import (
	"sales_monitor/scraper_app/shared/product/domain/entity"

	"github.com/playwright-community/playwright-go"
)

type Scraper func(context playwright.Browser, url string, wordsToIgnore []string) []*entity.ScrapedProduct

type ScraperConfig struct {
	ScrapingContent []ScrapingContent
	ScraperFunction Scraper
	MarketplaceName string
}

type ScrapingContent struct {
	URL           string
	Category      string
	WordsToIgnore []string
}
