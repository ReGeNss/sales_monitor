package mapper

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/feature/product/domain/entity"
)

func NotificationTaskToModel(e *entity.NotificationTask) *models.NotificationTask {
	if e == nil {
		return nil
	}
	products := make([]models.Product, 0, len(e.Products))
	for _, p := range e.Products {
		m := ProductToModel(p)
		if m != nil {
			products = append(products, *m)
		}
	}
	return &models.NotificationTask{
		BrandID:   e.BrandID,
		BrandName: e.BrandName,
		Products:  products,
	}
}
