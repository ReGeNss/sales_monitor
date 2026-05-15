package repository

import (
	"sales_monitor/scraper_app/feature/product/domain/entity"
	"sales_monitor/scraper_app/feature/product/domain/exception"
)

type PriceRepository interface {
	GetLatestProductPrice(productID int) (*entity.Price, exception.IDomainError)
}
