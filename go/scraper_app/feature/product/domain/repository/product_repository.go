package repository

import (
	"sales_monitor/scraper_app/feature/product/domain/entity"
	"sales_monitor/scraper_app/feature/product/domain/exception"
)

type ProductRepository interface {
	GetProductByFingerprint(fingerprint *string, brandID int, categoryID int, attributes []*entity.ProductAttribute) (*entity.Product, exception.IDomainError)
	FindSimilarCandidates(fingerprint *string, attributes []*entity.ProductAttribute, brandID int, categoryID int) ([]*entity.Product, exception.IDomainError)
	CreateProduct(product *entity.Product, attributes []*entity.ProductAttribute) (uint, exception.IDomainError)
	CreateProductAttribute(attribute *entity.ProductAttribute) exception.IDomainError
}
