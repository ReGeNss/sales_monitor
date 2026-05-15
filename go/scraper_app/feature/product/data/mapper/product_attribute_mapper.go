package mapper

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/feature/product/domain/entity"
)

func ProductAttributeToEntity(m *models.ProductAttribute) *entity.ProductAttribute {
	if m == nil {
		return nil
	}
	return &entity.ProductAttribute{
		ID:    m.AttributeID,
		Type:  m.AttributeType,
		Value: m.Value,
	}
}

func ProductAttributeToModel(e *entity.ProductAttribute) *models.ProductAttribute {
	if e == nil {
		return nil
	}
	return &models.ProductAttribute{
		AttributeID:   e.ID,
		AttributeType: e.Type,
		Value:         e.Value,
	}
}

func ProductAttributesToEntities(ms []models.ProductAttribute) []*entity.ProductAttribute {
	result := make([]*entity.ProductAttribute, 0, len(ms))
	for i := range ms {
		result = append(result, ProductAttributeToEntity(&ms[i]))
	}
	return result
}

func ProductAttributesToModels(es []*entity.ProductAttribute) []*models.ProductAttribute {
	result := make([]*models.ProductAttribute, 0, len(es))
	for _, e := range es {
		result = append(result, ProductAttributeToModel(e))
	}
	return result
}
