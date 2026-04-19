package main

import (
	"log"
	"os"

	"sales_monitor/internal/db"
	scraper_metrics "sales_monitor/scraper_app/feature/scraper/data/metrics"
	cached_scraped_product_repository "sales_monitor/scraper_app/feature/scraper/data/repository"
	scraper_storage "sales_monitor/scraper_app/feature/scraper/data/storage"
	scraper "sales_monitor/scraper_app/feature/scraper/domain/entity"
	cached_scraped_product_service "sales_monitor/scraper_app/feature/scraper/service"
	scraper_service "sales_monitor/scraper_app/feature/scraper/service"
	product_gateway "sales_monitor/scraper_app/shared/product/data/gateway"
	"sales_monitor/scraper_app/shared/product/data/repository"
	domainservice "sales_monitor/scraper_app/shared/product/domain/service"
	"sales_monitor/scraper_app/shared/product/service"
	"sales_monitor/scraper_app/shared/product/service/usecase"
)

func Run(plan scraper.ScrapingPlan) error {
	db.ConnectToDB()

	cachedScrapedProductService := cached_scraped_product_service.NewCachedScrapedProductService(
		cached_scraped_product_repository.NewCachedScrapedProductsRepository(db.GetDB()),
	)

	gormDB := db.GetDB()
	productRepo := repository.NewProductRepository(gormDB)
	categoryRepo := repository.NewCategoryRepository(gormDB)
	brandRepo := repository.NewBrandRepository(gormDB)
	marketplaceRepo := repository.NewMarketplaceRepository(gormDB)
	priceRepo := repository.NewPriceRepository(gormDB)
	matcher := domainservice.NewProductMatcher()

	productService := service.NewProductService(
		usecase.NewResolveCategoryUseCase(categoryRepo),
		usecase.NewResolveMarketplaceUseCase(marketplaceRepo),
		usecase.NewAssignBrandsUseCase(brandRepo),
		usecase.NewResolveBrandUseCase(brandRepo),
		usecase.NewResolveProductUseCase(productRepo, matcher),
		usecase.NewRecordPriceUseCase(marketplaceRepo, priceRepo),
		marketplaceRepo,
		product_gateway.NewNotificationPublisher(db.GetRedis()),
	)

	scraperService := scraper_service.NewScraperService(
		plan,
		productService,
		cachedScrapedProductService,
		scraper_storage.NewFileResultStorage(os.Getenv("SCRAPED_DATA_FOLDER")),
		scraper_metrics.NewPrometheusPublisher(),
	)

	scrapedProducts, err := scraperService.Scrape()
	if err != nil {
		return err
	}

	productService.ProcessProducts(scrapedProducts)
	log.Printf("processed %d categories", len(scrapedProducts))
	return nil
}
