package utils

import (
	scraper_config "sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/entity"
)

func CheckForProductUpdate(cachedProduct *scraper_config.LaterScrapedProductPrices, product *entity.ScrapedProduct) bool {
	res := cachedProduct.CurrentPrice == product.RegularPrice && (cachedProduct.DiscountedPrice == nil || *cachedProduct.DiscountedPrice == product.DiscountedPrice)
	return res
}