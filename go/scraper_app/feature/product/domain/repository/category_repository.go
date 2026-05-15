package repository

import (
	"sales_monitor/scraper_app/feature/product/domain/entity"
	"sales_monitor/scraper_app/feature/product/domain/exception"
)

type CategoryRepository interface {
	GetCategoryByName(name string) (*entity.Category, exception.IDomainError)
	CreateCategory(category *entity.Category) (uint, exception.IDomainError)
}
