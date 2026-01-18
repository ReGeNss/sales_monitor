package service

import (
	"sales_monitor/scraper_app/feature/scraper/domain/repository"
	scraper_config "sales_monitor/scraper_app/feature/scraper/domain/entity"
)


type CachedScrapedProductService interface {
	GetCachedScrapedProducts(marketplace string, category string) ([]*scraper_config.LaterScrapedProducts, error)
}

type cachedScrapedProductServiceImpl struct {
	repository repository.CachedScrapedProductsRepository
}

func NewCachedScrapedProductService(productRepository repository.CachedScrapedProductsRepository) CachedScrapedProductService {
	return &cachedScrapedProductServiceImpl{
		repository: productRepository,
	}
}

func (s *cachedScrapedProductServiceImpl) GetCachedScrapedProducts(marketplace string, category string) ([]*scraper_config.LaterScrapedProducts, error) {
	return s.repository.GetCachedScrapedProducts(marketplace, category)
}
