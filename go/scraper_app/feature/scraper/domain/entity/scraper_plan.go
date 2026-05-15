package entity

import "sales_monitor/scraper_app/feature/product/domain/entity"

type ScrapingPlan struct {
	Categories []ScrapingCategory
}

type ScrapingCategory struct {
	Category                     string
	ScrapersConfigs              []ScraperConfig
	WordsToIgnore                []string
	ProductDifferentiationEntity *entity.ProductDifferentiationEntity
}

type ScraperConfig struct {
	URLs   []string
	ShopID string
}

type LaterScrapedProducts map[string]LaterScrapedProductPrices

type LaterScrapedProductPrices struct {
	RegularPrice float64
	SpecialPrice *float64
}
