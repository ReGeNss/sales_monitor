package mapper

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/feature/product/domain/entity"
)

func PriceToEntity(m *models.Price) *entity.Price {
	if m == nil {
		return nil
	}
	return &entity.Price{
		ID:                   m.PriceID,
		MarketplaceProductID: m.MarketplaceProductID,
		RegularPrice:         m.RegularPrice,
		SpecialPrice:         m.SpecialPrice,
		CreatedAt:            m.CreatedAt,
	}
}
