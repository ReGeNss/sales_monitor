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

func (c *cachedScrapedProductsRepositoryImpl) GetCachedScrapedProducts(marketplace string, category string) (*entity.LaterScrapedProducts, error) {
	var prices []models.Price

	err := c.db.Transaction(func(tx *gorm.DB) error {
		var marketplaceModel models.Marketplace
		err := c.db.Model(&models.Marketplace{}).Where("name = ?", marketplace).First(&marketplaceModel).Error
		if err != nil {
			return err
		}

		var categoryModel models.Category
		err = c.db.Model(&models.Category{}).Where("name = ?", category).First(&categoryModel).Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Price{}).Table("prices as prc").Joins(
			"JOIN products p ON p.product_id = prc.product_id",
		).Where("marketplace_id = ? AND category_id = ?", marketplaceModel.MarketplaceID, categoryModel.CategoryID).Find(&prices).Error
		return err
	})

	if err != nil {
		return nil, err
	}

	laterScrapedProductsMap := make(map[string]entity.LaterScrapedProductPrices)
	for _, price := range prices {
		laterScrapedProductsMap[price.MarketplaceProduct.URL] = entity.LaterScrapedProductPrices{
			CurrentPrice:    price.RegularPrice,
			DiscountedPrice: price.DiscountPrice,
		}
	}

	laterScrapedProducts := entity.LaterScrapedProducts(laterScrapedProductsMap)
	return &laterScrapedProducts, nil
}
