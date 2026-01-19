package repository

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/shared/product/domain/entity"
)

type ProductRepository interface {
	GetProductByFingerprint(fingerprint *string, brandID int, categoryID int, attributes []*models.ProductAttribute) (*models.Product, error)
	GetMostSimilarProductID(fingerprint *string, attributes []*models.ProductAttribute, productDifferentiationEntity *entity.ProductDifferentiationEntity, brandID int, categoryID int, currentMarketplaceID int) (uint, error)
	CreateProduct(product *models.Product, attributes []*models.ProductAttribute) (uint,error)
	AddPriceToProduct(price *models.Price) error

	CreateCategory(category *models.Category) (uint, error)
	CreateBrand(brand *models.Brand) (uint, error)
	CreateMarketplace(marketplace *models.Marketplace) (uint, error)
	CreateProductAttribute(attribute *models.ProductAttribute) error

	GetCategoryByName(name string) (*models.Category, error)
	GetBrandByName(name string) (*models.Brand, error)
	GetMarketplaceByName(name string) (*models.Marketplace, error)
	GetAllBrands() ([]models.Brand, error)
	GetLaterScrapedProducts(brandID int) (entity.LaterScrapedProductsUrls, error)
}
