package event

import "sales_monitor/scraper_app/feature/product/domain/entity"

func NewDroppedProduct(p *entity.Product) DroppedProduct {
	return DroppedProduct{
		ID:       p.ID,
		Name:     p.Name,
		ImageURL: p.ImageURL.Url(),
	}
}

func NewDroppedProducts(products []*entity.Product) []DroppedProduct {
	result := make([]DroppedProduct, 0, len(products))
	for _, p := range products {
		result = append(result, NewDroppedProduct(p))
	}
	return result
}
