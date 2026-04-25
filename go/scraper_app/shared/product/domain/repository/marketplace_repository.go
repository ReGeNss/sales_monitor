package repository

import (
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/exception"
)

type MarketplaceRepository interface {
	GetMarketplaceByName(name string) (*entity.Marketplace, exception.IDomainError)
	CreateMarketplace(marketplace *entity.Marketplace) (uint, exception.IDomainError)
	GetLaterScrapedProducts(brandID int) (entity.LaterScrapedProductsUrls, exception.IDomainError)
	AddPriceToMarketplaceProduct(productID int, marketplaceID int, url string, regularPrice float64, specialPrice *float64) exception.IDomainError
	AddPriceToMarketplaceProductID(marketplaceProductID int, regularPrice float64, specialPrice *float64) exception.IDomainError
}
