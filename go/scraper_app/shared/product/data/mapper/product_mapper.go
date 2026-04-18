package mapper

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/shared/product/domain/entity"
)

func ProductToEntity(m *models.Product) *entity.Product {
	if m == nil {
		return nil
	}
	return &entity.Product{
		ID:              m.ProductID,
		Name:            m.Name,
		NameFingerprint: m.NameFingerprint,
		ImageURL:        m.ImageURL,
		BrandID:         m.BrandID,
		CategoryID:      m.CategoryID,
		Attributes:      ProductAttributesToEntities(m.Attributes),
	}
}

func ProductToModel(e *entity.Product) *models.Product {
	if e == nil {
		return nil
	}
	return &models.Product{
		ProductID:       e.ID,
		Name:            e.Name,
		NameFingerprint: e.NameFingerprint,
		ImageURL:        e.ImageURL,
		BrandID:         e.BrandID,
		CategoryID:      e.CategoryID,
	}
}
