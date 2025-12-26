package service

import "sales_monitor/scraper/shared/product/domain/entity"

type ProductService interface {
	ProcessProducts(products []entity.Product)
}
