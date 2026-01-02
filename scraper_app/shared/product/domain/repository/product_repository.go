package repository

import (
	"sales_monitor/internal/models"
)

type ProductRepository interface {
	GetProductByFingerprint(fingerprint string, brandID int, categoryID int, attributes []*models.ProductAttribute) (*models.Product, error)
	GetMostSimilarProductID(fingerprint string, attributes []*models.ProductAttribute, brandID int, categoryID int) (uint, error)
	CreateProduct(product *models.Product, attributes []*models.ProductAttribute) (uint,error)
	AddPriceToProduct(price *models.Price) error

	CreateCategory(category *models.Category) (uint, error)
	CreateBrand(brand *models.Brand) (uint, error)
	CreateMarketplace(marketplace *models.Marketplace) (uint, error)
	CreateProductAttribute(attribute *models.ProductAttribute) error

	GetCategoryByName(name string) (*models.Category, error)
	GetBrandByName(name string) (*models.Brand, error)
	GetMarketplaceByName(name string) (*models.Marketplace, error)
}
