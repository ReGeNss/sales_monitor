package scraper

import (
	"github.com/playwright-community/playwright-go"
	"sales_monitor/scraper_app/shared/product/domain/entity"
)

type Scraper func(context playwright.Browser, url string) []*entity.Product

type ScraperConfig struct {
	URLs []string
	Scraper Scraper
}