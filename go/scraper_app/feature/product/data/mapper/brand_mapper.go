package mapper

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/feature/product/domain/entity"
)

func BrandToEntity(m *models.Brand) *entity.Brand {
	if m == nil {
		return nil
	}
	return &entity.Brand{
		ID:        m.BrandID,
		Name:      m.Name,
		BannerURL: m.BannerURL,
	}
}

func BrandToModel(e *entity.Brand) *models.Brand {
	if e == nil {
		return nil
	}
	return &models.Brand{
		BrandID:   e.ID,
		Name:      e.Name,
		BannerURL: e.BannerURL,
	}
}

func BrandsToEntities(ms []models.Brand) []*entity.Brand {
	result := make([]*entity.Brand, 0, len(ms))
	for i := range ms {
		result = append(result, BrandToEntity(&ms[i]))
	}
	return result
}
