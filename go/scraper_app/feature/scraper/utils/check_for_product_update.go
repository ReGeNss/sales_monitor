package utils

import (
	scraper_config "sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/service/dto"
)

func CheckForProductUpdate(cachedProduct *scraper_config.LaterScrapedProductPrices, product *dto.ScrapedProductDto) bool {
	res := cachedProduct.RegularPrice == product.RegularPrice && (cachedProduct.SpecialPrice == nil || *cachedProduct.SpecialPrice == product.SpecialPrice)
	return res
}