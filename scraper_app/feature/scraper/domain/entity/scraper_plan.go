package entity

import "sales_monitor/scraper_app/shared/product/domain/entity"

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
	URLs    []string
	Scraper Scraper
}

type LaterScrapedProducts map[string]LaterScrapedProductPrices

type LaterScrapedProductPrices struct {
	CurrentPrice    float64
	DiscountedPrice *float64
}
