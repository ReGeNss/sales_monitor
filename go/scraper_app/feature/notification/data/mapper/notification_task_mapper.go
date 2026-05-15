package mapper

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/feature/notification/domain/entity"
)

func NotificationTaskToModel(e *entity.NotificationTask) *models.NotificationTask {
	if e == nil {
		return nil
	}
	products := make([]models.NotificationProduct, 0, len(e.Products))
	for _, p := range e.Products {
		m := &models.NotificationProduct{
			ID:       p.ID,
			Name:     p.Name,
			ImageURL: p.ImageURL,
		}
		
			products = append(products, *m)
		
	}
	return &models.NotificationTask{
		BrandID:   e.BrandID,
		BrandName: e.BrandName,
		Products:  products,
	}
}
