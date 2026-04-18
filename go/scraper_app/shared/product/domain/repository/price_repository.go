package repository

import "sales_monitor/scraper_app/shared/product/domain/entity"

type PriceRepository interface {
	GetLatestProductPrice(productID int) (*entity.Price, error)
}
