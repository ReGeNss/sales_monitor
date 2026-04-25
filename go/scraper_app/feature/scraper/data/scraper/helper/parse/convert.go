package parse

import (
	"sales_monitor/scraper_app/feature/scraper/data/scraper/model"
	entity "sales_monitor/scraper_app/shared/product/domain/entity"
)

func ToEntityProducts(raw []*model.ScrapedProduct, wordsToIgnore []string) []*entity.ScrapedProduct {
	products := make([]*entity.ScrapedProduct, 0, len(raw))
	for _, p := range raw {
		product, _ := entity.NewScrapedProduct(
			p.Name,
			p.RegularPrice,
			p.SpecialPrice,
			p.ImageURL,
			p.URL,
			p.BrandName,
			p.Volume,
			p.Weight,
			wordsToIgnore,
		)
		if product != nil {
			products = append(products, product)
		}
	}
	return products
}
