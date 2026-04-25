package mapper

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/feature/scraper/domain/exception"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	valueObject "sales_monitor/scraper_app/shared/product/domain/entity/value_object"
)

func ProductToEntity(m *models.Product) (*entity.Product, exception.IDomainError) {
	if m == nil {
		return nil, exception.NewDomainError("Product is nil")
	}

	imageUrl, err := valueObject.NewUrl(m.ImageURL) 

	if err != nil {
		return nil, err
	}

	return &entity.Product{
		ID:              m.ProductID,
		Name:            m.Name,
		NameFingerprint: m.NameFingerprint,
		ImageURL:        *imageUrl,
		BrandID:         m.BrandID,
		CategoryID:      m.CategoryID,
		Attributes:      ProductAttributesToEntities(m.Attributes),
	}, nil
}

func ProductToModel(e *entity.Product) *models.Product {
	if e == nil {
		return nil
	}
	return &models.Product{
		ProductID:       e.ID,
		Name:            e.Name,
		NameFingerprint: e.NameFingerprint,
		ImageURL:        e.ImageURL.Url(),
		BrandID:         e.BrandID,
		CategoryID:      e.CategoryID,
	}
}
