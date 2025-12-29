package repository

import (
	"sales_monitor/internal/models"
)

type ProductRepository interface {
	GetProductByFingerprint(fingerprint string) (*models.Product, error)
	GetMostSimilarProductID(fingerprint string) (uint, error)
	CreateProduct(product *models.Product) (uint,error)
	AddPriceToProduct(price *models.Price) error

	CreateCategory(category *models.Category) (uint, error)
	CreateBrand(brand *models.Brand) (uint, error)
	CreateMarketplace(marketplace *models.Marketplace) (uint, error)

	GetCategoryByName(name string) (*models.Category, error)
	GetBrandByName(name string) (*models.Brand, error)
	GetMarketplaceByName(name string) (*models.Marketplace, error)
}
