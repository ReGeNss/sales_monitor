package model

import "sales_monitor/scraper_app/shared/product/domain/entity"

type ScrapeResult struct {
	Products   []*entity.Product
	FoundCount int
	NewCount   int
}
