package main

import (
	"log"
	"os"

	"sales_monitor/internal/db"
	notification_repository "sales_monitor/scraper_app/feature/notification/data/repository"
	notificationEntity "sales_monitor/scraper_app/feature/notification/domain/entity"
	"sales_monitor/scraper_app/feature/product/data/repository"
	productevent "sales_monitor/scraper_app/feature/product/domain/event"
	domainservice "sales_monitor/scraper_app/feature/product/domain/service"
	"sales_monitor/scraper_app/feature/product/service"
	"sales_monitor/scraper_app/feature/product/service/usecase"
	scraper_logging "sales_monitor/scraper_app/feature/scraper/data/logging"
	scraper_metrics "sales_monitor/scraper_app/feature/scraper/data/metrics"
	cached_scraped_product_repository "sales_monitor/scraper_app/feature/scraper/data/repository"
	scraper_factory "sales_monitor/scraper_app/feature/scraper/data/scraper"
	scraper_storage "sales_monitor/scraper_app/feature/scraper/data/storage"
	scraper "sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/domain/gateway"
	scraper_service "sales_monitor/scraper_app/feature/scraper/service"
	"sales_monitor/scraper_app/utils"
)

func Run(plan scraper.ScrapingPlan) error {
	db.ConnectToDB()

	cachedScrapedProductService :=
		cached_scraped_product_repository.NewCachedScrapedProductsRepository(db.GetDB())

	gormDB := db.GetDB()
	productRepo := repository.NewProductRepository(gormDB)
	categoryRepo := repository.NewCategoryRepository(gormDB)
	brandRepo := repository.NewBrandRepository(gormDB)
	marketplaceRepo := repository.NewMarketplaceRepository(gormDB)
	priceRepo := repository.NewPriceRepository(gormDB)
	matcher := domainservice.NewProductMatcher()

	productServiceEventBus := utils.NewEventBus()

	notificationRepository := notification_repository.NewNotificationRepository(db.GetRedis())

	productServiceEventBus.Subscribe(&productevent.PriceDropDetected{}, func(payload interface{}) {
		e, ok := payload.(*productevent.PriceDropDetected)
		if !ok {
			log.Printf("unexpected event type for PriceDropDetected handler: %T", payload)
			return
		}
		products := make([]*notificationEntity.Product, 0, len(e.Products))
		for _, p := range e.Products {
			products = append(products, &notificationEntity.Product{
				ID:       p.ID,
				Name:     p.Name,
				ImageURL: p.ImageURL,
			})
		}
		task := &notificationEntity.NotificationTask{
			BrandID:   e.BrandID,
			BrandName: e.BrandName,
			Products:  products,
		}
		if err := notificationRepository.SendNotification(task); err != nil {
			log.Printf("could not send notification: %v", err)
		}
	})

	productService := service.NewProductService(
		usecase.NewResolveCategoryUseCase(categoryRepo),
		usecase.NewResolveMarketplaceUseCase(marketplaceRepo),
		usecase.NewAssignBrandsUseCase(brandRepo),
		usecase.NewResolveBrandUseCase(brandRepo),
		usecase.NewResolveProductUseCase(productRepo, matcher),
		usecase.NewRecordPriceUseCase(marketplaceRepo, priceRepo),
		marketplaceRepo,
		productServiceEventBus,
	)

	errorLogger := scraper_logging.NewScreenshotErrorLogger()

	scraperFactory, err := scraper_factory.NewScraperFactory(errorLogger)
	if err != nil {
		return err
	}
	defer scraperFactory.Close()

	scraperServiceEventBus := utils.NewEventBus()

	metricsPublisher := scraper_metrics.NewPrometheusPublisher()
	scraperServiceEventBus.Subscribe(&scraper.ScrapingCompleted{}, func(payload interface{}) {
		event, ok := payload.(*scraper.ScrapingCompleted)
		if !ok {
			log.Printf("unexpected event type for ScrapingCompletedEvent handler: %T", payload)
			return
		}
		metricsPublisher.Publish(gateway.ScrapingMetrics{
			Found:   event.Found,
			Scraped: event.Scraped,
			New:     event.New,
			OnSale:  event.OnSale,
		}, event.Results)
	})

	scraperService := scraper_service.NewScraperService(
		plan,
		cachedScrapedProductService,
		scraperFactory,
		scraperServiceEventBus,
	)

	resultStorage := scraper_storage.NewFileResultStorage(os.Getenv("SCRAPED_DATA_FOLDER"))

	scrapedProducts, err := scraperService.Scrape()
	if err != nil {
		return err
	}

	names := make([]string, 0, len(plan.Categories))
	for _, c := range plan.Categories {
		names = append(names, c.Category)
	}

	resultStorage.Save(scrapedProducts, names)
	productService.ProcessProducts(scrapedProducts)
	log.Printf("processed %d categories", len(scrapedProducts))
	return nil
}
