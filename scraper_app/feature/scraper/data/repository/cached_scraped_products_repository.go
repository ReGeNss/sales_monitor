package repository

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/domain/repository"

	"gorm.io/gorm"
)

type cachedScrapedProductsRepositoryImpl struct {
	db *gorm.DB
}


func NewCachedScrapedProductsRepository(db *gorm.DB) repository.CachedScrapedProductsRepository {
	return &cachedScrapedProductsRepositoryImpl{
		db: db,
	}
}

func (c *cachedScrapedProductsRepositoryImpl) GetCachedScrapedProducts(marketplace string, category string) ([]*entity.LaterScrapedProducts, error) {
	var laterScrapedProducts []*entity.LaterScrapedProducts

	err := c.db.Transaction(func(tx *gorm.DB) error {
		var marketplace models.Marketplace
		err := c.db.Model(&models.Marketplace{}).Where("name = ?", marketplace).First(&marketplace).Error
		if err != nil {
			return err
		}

		var category models.Category
		err = c.db.Model(&models.Category{}).Where("name = ?", category).First(&category).Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Price{}).Where("marketplace_id = ? AND category_id = ?", marketplace.MarketplaceID, category.CategoryID).Find(&laterScrapedProducts).Error
		return err
	})

	if err != nil {
		return nil, err
	}

	return laterScrapedProducts, nil
}
