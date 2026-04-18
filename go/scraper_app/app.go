package main

import (
	"log"

	"sales_monitor/internal/db"
	"sales_monitor/scraper_app/core/api"
	cached_scraped_product_repository "sales_monitor/scraper_app/feature/scraper/data/repository"
	scraper "sales_monitor/scraper_app/feature/scraper/domain/entity"
	cached_scraped_product_service "sales_monitor/scraper_app/feature/scraper/service"
	scraper_service "sales_monitor/scraper_app/feature/scraper/service"
	"sales_monitor/scraper_app/shared/product/data/repository"
	"sales_monitor/scraper_app/shared/product/service"
)

func Run(plan scraper.ScrapingPlan) error {
	db.ConnectToDB()

	cachedScrapedProductService := cached_scraped_product_service.NewCachedScrapedProductService(
		cached_scraped_product_repository.NewCachedScrapedProductsRepository(db.GetDB()),
	)

	productService := service.NewProductService(
		repository.NewProductRepository(db.GetDB(), api.NewHTTPClient(), db.GetRedis()),
	)

	scraperService := scraper_service.NewScraperService(
		plan,
		productService,
		cachedScrapedProductService,
	)

	scrapedProducts, err := scraperService.Scrape()
	if err != nil {
		return err
	}

	productService.ProcessProducts(scrapedProducts)
	log.Printf("processed %d categories", len(scrapedProducts))
	return nil
}
