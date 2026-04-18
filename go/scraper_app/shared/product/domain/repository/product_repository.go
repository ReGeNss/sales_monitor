package repository

import (
	"sales_monitor/scraper_app/shared/product/domain/entity"
)

type ProductRepository interface {
	GetProductByFingerprint(fingerprint *string, brandID int, categoryID int, attributes []*entity.ProductAttribute) (*entity.Product, error)
	FindSimilarCandidates(fingerprint *string, attributes []*entity.ProductAttribute, brandID int, categoryID int) ([]*entity.Product, error)
	CreateProduct(product *entity.Product, attributes []*entity.ProductAttribute) (uint, error)
	AddPriceToMarketplaceProduct(productID int, marketplaceID int, url string, regularPrice float64, specialPrice *float64) error
	AddPriceToMarketplaceProductID(marketplaceProductID int, regularPrice float64, specialPrice *float64) error

	CreateCategory(category *entity.Category) (uint, error)
	CreateBrand(brand *entity.Brand) (uint, error)
	CreateMarketplace(marketplace *entity.Marketplace) (uint, error)
	CreateProductAttribute(attribute *entity.ProductAttribute) error

	GetCategoryByName(name string) (*entity.Category, error)
	GetBrandByName(name string) (*entity.Brand, error)
	GetMarketplaceByName(name string) (*entity.Marketplace, error)
	GetAllBrands() ([]*entity.Brand, error)
	GetLaterScrapedProducts(brandID int) (entity.LaterScrapedProductsUrls, error)
	GetLatestProductPrice(productID int) (*entity.Price, error)

	SendNotification(notificationTask *entity.NotificationTask) error
}
