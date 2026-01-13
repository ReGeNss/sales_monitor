package entity

import (
	"sales_monitor/scraper_app/shared/product/domain/entity"

	"github.com/playwright-community/playwright-go"
)

type Scraper func(context playwright.Browser, url string, wordsToIgnore []string) []*entity.ScrapedProduct
