package repository

import (
	"fmt"
	database "sales_monitor/internal/db"
	"sales_monitor/internal/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProductByFingerprint(fingerprint string) (*models.Product, error)
	GetMostSimilarProductID(fingerprint string) (uint, error)
	CreateProduct(product *models.Product) error
	AddPriceToProduct(price *models.Price) error
}

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepositoryImpl{
		db: database.GetDB(),
	}
}

func (p *productRepositoryImpl) AddPriceToProduct(price *models.Price) error {
	return p.db.Model(&models.Price{}).Create(price).Error
}

func (p *productRepositoryImpl) CreateProduct(product *models.Product) error {
	return p.db.Model(&models.Product{}).Create(product).Error
}

func (p *productRepositoryImpl) GetMostSimilarProductID(fingerprint string) (uint, error) {
	var product models.Product
	err := p.db.Model(&models.Product{}).
		Select("product_id").
		Where("MATCH(name_fingerprint) AGAINST(? IN NATURAL LANGUAGE MODE) > 0.9", fingerprint).
		Order(fmt.Sprintf("MATCH(name_fingerprint) AGAINST('%s' IN NATURAL LANGUAGE MODE) DESC", fingerprint)).
		Limit(1).
		First(&product).Error
	if err != nil {
		return 0, err
	}
	return uint(product.ProductID), nil
}

func (p *productRepositoryImpl) GetProductByFingerprint(fingerprint string) (*models.Product, error) {
	var product models.Product
	err := p.db.Model(&models.Product{}).Where("name_fingerprint = ?", fingerprint).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

