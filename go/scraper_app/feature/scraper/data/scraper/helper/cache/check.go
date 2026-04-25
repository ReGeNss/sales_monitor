package cache

import "sales_monitor/scraper_app/shared/product/domain/entity"
import later "sales_monitor/scraper_app/feature/scraper/domain/entity"

func IsUpToDate(cached *later.LaterScrapedProductPrices, product *entity.ScrapedProduct) bool {
	return cached.RegularPrice == product.RegularPrice() &&
		(cached.SpecialPrice == nil || *cached.SpecialPrice == product.SpecialPrice())
}
