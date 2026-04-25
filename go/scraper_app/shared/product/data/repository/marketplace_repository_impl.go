package repository

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/shared/product/data/mapper"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/exception"
	"sales_monitor/scraper_app/shared/product/domain/repository"

	"gorm.io/gorm"
)

type marketplaceRepositoryImpl struct {
	db *gorm.DB
}

func NewMarketplaceRepository(db *gorm.DB) repository.MarketplaceRepository {
	return &marketplaceRepositoryImpl{db: db}
}

func (r *marketplaceRepositoryImpl) GetMarketplaceByName(name string) (*entity.Marketplace, exception.IDomainError) {
	var marketplace models.Marketplace
	err := r.db.Model(&models.Marketplace{}).Where("name = ?", name).First(&marketplace).Error
	if err != nil {
		return nil, err
	}
	return mapper.MarketplaceToEntity(&marketplace), nil
}

func (r *marketplaceRepositoryImpl) CreateMarketplace(marketplace *entity.Marketplace) (uint, exception.IDomainError) {
	m := mapper.MarketplaceToModel(marketplace)
	if err := r.db.Model(&models.Marketplace{}).Create(m).Error; err != nil {
		return 0, err
	}
	marketplace.ID = m.MarketplaceID
	return uint(m.MarketplaceID), nil
}

func (r *marketplaceRepositoryImpl) GetLaterScrapedProducts(brandID int) (entity.LaterScrapedProductsUrls, exception.IDomainError) {
	var marketplaceProducts []models.MarketplaceProduct
	err := r.db.Model(&models.MarketplaceProduct{}).Table("marketplace_products as mp").
		Joins("JOIN products p ON p.product_id = mp.product_id").
		Where("p.brand_id = ?", brandID).
		Find(&marketplaceProducts).Error
	if err != nil {
		return nil, err
	}

	urls := make(entity.LaterScrapedProductsUrls)
	for _, mp := range marketplaceProducts {
		urls[mp.URL] = mp.MarketplaceProductID
	}
	return urls, nil
}

func (r *marketplaceRepositoryImpl) AddPriceToMarketplaceProduct(productID int, marketplaceID int, url string, regularPrice float64, specialPrice *float64) exception.IDomainError {
	marketplaceProduct := models.MarketplaceProduct{
		MarketplaceID: marketplaceID,
		ProductID:     productID,
		URL:           url,
	}

	if err := r.db.Model(&models.MarketplaceProduct{}).
		Where("marketplace_id = ? AND product_id = ? AND url = ?", marketplaceID, productID, url).
		FirstOrCreate(&marketplaceProduct).Error; err != nil {
		return err
	}

	price := models.Price{
		MarketplaceProductID: marketplaceProduct.MarketplaceProductID,
		RegularPrice:         regularPrice,
		SpecialPrice:         specialPrice,
	}

	return r.db.Model(&models.Price{}).Create(&price).Error
}

func (r *marketplaceRepositoryImpl) AddPriceToMarketplaceProductID(marketplaceProductID int, regularPrice float64, specialPrice *float64) exception.IDomainError {
	price := models.Price{
		MarketplaceProductID: marketplaceProductID,
		RegularPrice:         regularPrice,
		SpecialPrice:         specialPrice,
	}

	return r.db.Model(&models.Price{}).Create(&price).Error
}
