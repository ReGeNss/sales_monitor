package entity

import "sales_monitor/scraper_app/shared/product/domain/entity"

type ScrapingResult struct {
	ScrapedProducts              []*entity.ScrapedProducts
	ProductDifferentiationEntity *entity.ProductDifferentiationEntity
}
