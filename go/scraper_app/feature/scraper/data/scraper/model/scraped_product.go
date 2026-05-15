package model

import "sales_monitor/scraper_app/feature/product/domain/entity"

type ScrapeResult struct {
	Products   []*entity.Product
	FoundCount int
	NewCount   int
}
