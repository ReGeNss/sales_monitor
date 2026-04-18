package repository

import "sales_monitor/scraper_app/shared/product/domain/entity"

type MarketplaceRepository interface {
	GetMarketplaceByName(name string) (*entity.Marketplace, error)
	CreateMarketplace(marketplace *entity.Marketplace) (uint, error)
	GetLaterScrapedProducts(brandID int) (entity.LaterScrapedProductsUrls, error)
	AddPriceToMarketplaceProduct(productID int, marketplaceID int, url string, regularPrice float64, specialPrice *float64) error
	AddPriceToMarketplaceProductID(marketplaceProductID int, regularPrice float64, specialPrice *float64) error
}
