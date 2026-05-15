package entity

import "sales_monitor/scraper_app/feature/product/domain/entity"

type ScrapingResult struct {
	ScrapedProducts              []*entity.ScrapedProducts
	ProductDifferentiationEntity *entity.ProductDifferentiationEntity
}
