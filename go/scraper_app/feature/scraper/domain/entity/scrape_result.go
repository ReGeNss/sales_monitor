package entity

import "sales_monitor/scraper_app/feature/product/domain/entity"

type ScrapeResult struct {
	Products   []*entity.ScrapedProduct
	FoundCount int
	NewCount   int
}
