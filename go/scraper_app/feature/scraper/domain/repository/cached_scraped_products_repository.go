package repository

import (
	"sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/domain/exception"
)

type CachedScrapedProductsRepository interface {
	GetCachedScrapedProducts(marketplace string, category string) (*entity.LaterScrapedProducts, exception.IDomainError)
}
