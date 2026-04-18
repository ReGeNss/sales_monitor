package mapper

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/shared/product/domain/entity"
)

func MarketplaceToEntity(m *models.Marketplace) *entity.Marketplace {
	if m == nil {
		return nil
	}
	return &entity.Marketplace{
		ID:   m.MarketplaceID,
		Name: m.Name,
		URL:  m.URL,
	}
}

func MarketplaceToModel(e *entity.Marketplace) *models.Marketplace {
	if e == nil {
		return nil
	}
	return &models.Marketplace{
		MarketplaceID: e.ID,
		Name:          e.Name,
		URL:           e.URL,
	}
}
