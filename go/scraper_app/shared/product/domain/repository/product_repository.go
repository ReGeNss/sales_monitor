package repository

import "sales_monitor/scraper_app/shared/product/domain/entity"

type ProductRepository interface {
	GetProductByFingerprint(fingerprint *string, brandID int, categoryID int, attributes []*entity.ProductAttribute) (*entity.Product, error)
	FindSimilarCandidates(fingerprint *string, attributes []*entity.ProductAttribute, brandID int, categoryID int) ([]*entity.Product, error)
	CreateProduct(product *entity.Product, attributes []*entity.ProductAttribute) (uint, error)
	CreateProductAttribute(attribute *entity.ProductAttribute) error
}
