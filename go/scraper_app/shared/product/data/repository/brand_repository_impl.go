package repository

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/shared/product/data/mapper"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/repository"

	"gorm.io/gorm"
)

type brandRepositoryImpl struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) repository.BrandRepository {
	return &brandRepositoryImpl{db: db}
}

func (r *brandRepositoryImpl) GetBrandByName(name string) (*entity.Brand, error) {
	var brand models.Brand
	err := r.db.Model(&models.Brand{}).Where("name = ?", name).First(&brand).Error
	if err != nil {
		return nil, err
	}
	return mapper.BrandToEntity(&brand), nil
}

func (r *brandRepositoryImpl) GetAllBrands() ([]*entity.Brand, error) {
	var brands []models.Brand
	err := r.db.Model(&models.Brand{}).Find(&brands).Error
	if err != nil {
		return nil, err
	}
	return mapper.BrandsToEntities(brands), nil
}

func (r *brandRepositoryImpl) CreateBrand(brand *entity.Brand) (uint, error) {
	m := mapper.BrandToModel(brand)
	if err := r.db.Model(&models.Brand{}).Create(m).Error; err != nil {
		return 0, err
	}
	brand.ID = m.BrandID
	return uint(m.BrandID), nil
}
