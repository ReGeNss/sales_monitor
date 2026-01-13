package entity

import "sales_monitor/scraper_app/shared/product/domain/entity"

type ScrapingContent struct {
	URL           string
	Category      string
	WordsToIgnore []string
	ProductDifferentiationEntity *entity.ProductDifferentiationEntity
}

