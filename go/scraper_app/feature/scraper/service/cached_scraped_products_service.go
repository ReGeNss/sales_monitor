package service

import (
	scraper_config "sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/domain/exception"
	"sales_monitor/scraper_app/feature/scraper/domain/repository"
)

type CachedScrapedProductService interface {
	GetCachedScrapedProducts(marketplace string, category string) (*scraper_config.LaterScrapedProducts, exception.IDomainError)
}

type cachedScrapedProductServiceImpl struct {
	repository repository.CachedScrapedProductsRepository
}

func NewCachedScrapedProductService(productRepository repository.CachedScrapedProductsRepository) CachedScrapedProductService {
	return &cachedScrapedProductServiceImpl{
		repository: productRepository,
	}
}

func (s *cachedScrapedProductServiceImpl) GetCachedScrapedProducts(marketplace string, category string) (*scraper_config.LaterScrapedProducts, exception.IDomainError) {
	return s.repository.GetCachedScrapedProducts(marketplace, category)
}
