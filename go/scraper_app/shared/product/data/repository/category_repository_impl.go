package repository

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/shared/product/data/mapper"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/repository"

	"gorm.io/gorm"
)

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepositoryImpl{db: db}
}

func (r *categoryRepositoryImpl) GetCategoryByName(name string) (*entity.Category, error) {
	var category models.Category
	err := r.db.Model(&models.Category{}).Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return mapper.CategoryToEntity(&category), nil
}

func (r *categoryRepositoryImpl) CreateCategory(category *entity.Category) (uint, error) {
	m := mapper.CategoryToModel(category)
	if err := r.db.Model(&models.Category{}).Create(m).Error; err != nil {
		return 0, err
	}
	category.ID = m.CategoryID
	return uint(m.CategoryID), nil
}
