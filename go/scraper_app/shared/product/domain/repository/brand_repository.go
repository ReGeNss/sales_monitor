package repository

import "sales_monitor/scraper_app/shared/product/domain/entity"

type BrandRepository interface {
	GetBrandByName(name string) (*entity.Brand, error)
	GetAllBrands() ([]*entity.Brand, error)
	CreateBrand(brand *entity.Brand) (uint, error)
}
