package mapper

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/shared/product/domain/entity"
)

func CategoryToEntity(m *models.Category) *entity.Category {
	if m == nil {
		return nil
	}
	return &entity.Category{
		ID:   m.CategoryID,
		Name: m.Name,
	}
}

func CategoryToModel(e *entity.Category) *models.Category {
	if e == nil {
		return nil
	}
	return &models.Category{
		CategoryID: e.ID,
		Name:       e.Name,
	}
}
