package service

import "sales_monitor/scraper_app/shared/product/domain/entity"

type ProductService interface {
	ProcessProducts(products []entity.Product)
}
