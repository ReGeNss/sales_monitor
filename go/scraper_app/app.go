package main

import (
	"log"
	"os"

	"sales_monitor/internal/db"
	scraper_logging "sales_monitor/scraper_app/feature/scraper/data/logging"
	scraper_metrics "sales_monitor/scraper_app/feature/scraper/data/metrics"
	cached_scraped_product_repository "sales_monitor/scraper_app/feature/scraper/data/repository"
	scraper_storage "sales_monitor/scraper_app/feature/scraper/data/storage"
	scraper "sales_monitor/scraper_app/feature/scraper/domain/entity"
	scraper_service "sales_monitor/scraper_app/feature/scraper/service"
	product_gateway "sales_monitor/scraper_app/shared/product/data/gateway"
	"sales_monitor/scraper_app/shared/product/data/repository"
	domainservice "sales_monitor/scraper_app/shared/product/domain/service"
	"sales_monitor/scraper_app/shared/product/service"
	"sales_monitor/scraper_app/shared/product/service/usecase"
)

func Run(plan scraper.ScrapingPlan) error {
	db.ConnectToDB()

	cachedScrapedProductService := scraper_service.NewCachedScrapedProductService(
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
		scraper_logging.NewScreenshotErrorLogger(),
	)

	scrapedProducts, err := scraperService.Scrape()
	if err != nil {
		return err
	}

	productService.ProcessProducts(scrapedProducts)
	log.Printf("processed %d categories", len(scrapedProducts))
	return nil
}
