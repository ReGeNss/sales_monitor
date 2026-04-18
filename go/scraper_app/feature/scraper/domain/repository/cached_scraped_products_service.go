package repository

import "sales_monitor/scraper_app/feature/scraper/domain/entity"

type CachedScrapedProductsRepository interface {
	GetCachedScrapedProducts(marketplace string, category string) (*entity.LaterScrapedProducts, error)
}