package repository

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/shared/product/data/mapper"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/exception"
	"sales_monitor/scraper_app/shared/product/domain/repository"

	"gorm.io/gorm"
)

type priceRepositoryImpl struct {
	db *gorm.DB
}

func NewPriceRepository(db *gorm.DB) repository.PriceRepository {
	return &priceRepositoryImpl{db: db}
}

func (r *priceRepositoryImpl) GetLatestProductPrice(productID int) (*entity.Price, exception.IDomainError) {
	var price models.Price
	err := r.db.Model(&models.Price{}).
		Joins("JOIN marketplace_products mp ON mp.marketplace_product_id = prices.marketplace_product_id").
		Where("mp.product_id = ?", productID).
		Order("prices.created_at DESC, prices.price_id DESC").
		First(&price).Error
	if err != nil {
		return nil, err
	}
	return mapper.PriceToEntity(&price), nil
}
