package repository

import "sales_monitor/scraper_app/shared/product/domain/entity"

type CategoryRepository interface {
	GetCategoryByName(name string) (*entity.Category, error)
	CreateCategory(category *entity.Category) (uint, error)
}
