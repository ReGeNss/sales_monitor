package repository

import (
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/exception"
)

type CategoryRepository interface {
	GetCategoryByName(name string) (*entity.Category, exception.IDomainError)
	CreateCategory(category *entity.Category) (uint, exception.IDomainError)
}
