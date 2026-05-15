package repository

import (
	"sales_monitor/scraper_app/feature/product/domain/entity"
	"sales_monitor/scraper_app/feature/product/domain/exception"
)

type BrandRepository interface {
	GetBrandByName(name string) (*entity.Brand, exception.IDomainError)
	GetAllBrands() ([]*entity.Brand, exception.IDomainError)
	CreateBrand(brand *entity.Brand) (uint, exception.IDomainError)
}
