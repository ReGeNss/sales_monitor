package mapper

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/feature/product/domain/entity"
)

func MarketplaceProductToEntity(m *models.MarketplaceProduct) *entity.MarketplaceProduct {
	if m == nil {
		return nil
	}
	return &entity.MarketplaceProduct{
		ID:            m.MarketplaceProductID,
		MarketplaceID: m.MarketplaceID,
		ProductID:     m.ProductID,
		URL:           m.URL,
	}
}
