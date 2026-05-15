package main

import (
	"log"
	"os"

	"sales_monitor/internal/db"
	"sales_monitor/scraper_app/feature/notification"
	notificationEntity "sales_monitor/scraper_app/feature/notification/domain/entity"
	"sales_monitor/scraper_app/feature/product"
	productEvent "sales_monitor/scraper_app/feature/product/domain/event"
	"sales_monitor/scraper_app/feature/scraper"
	scraperEntity "sales_monitor/scraper_app/feature/scraper/domain/entity"
	scraperEvent "sales_monitor/scraper_app/feature/scraper/domain/event"
	"sales_monitor/scraper_app/feature/statistics"
	statisticsEntity "sales_monitor/scraper_app/feature/statistics/domain/entity"
	"sales_monitor/scraper_app/feature/storage"
	"sales_monitor/scraper_app/utils"
)

func Run(plan scraperEntity.ScrapingPlan) error {
	db.ConnectToDB()
	gormDB := db.GetDB()

	productEventBus := utils.NewEventBus()
	scraperEventBus := utils.NewEventBus()

	notificationRepository := notification.NewNotificationFeature(db.GetRedis())
	statisticsRepository := statistics.NewStatisticsRepository()
	resultStorage := storage.NewResultStorage(os.Getenv("SCRAPED_DATA_FOLDER"))

	productService := product.NewProductService(gormDB, productEventBus)

	scraperService, scraperFactory, err := scraper.NewScraperService(plan, gormDB, scraperEventBus)
	if err != nil {
		return err
	}
	defer scraperFactory.Close()

	productEventBus.Subscribe(&productEvent.PriceDropDetected{}, func(payload interface{}) {
		e, ok := payload.(*productEvent.PriceDropDetected)
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

	scraperEventBus.Subscribe(&scraperEvent.ScrapingCompleted{}, func(payload interface{}) {
		e, ok := payload.(*scraperEvent.ScrapingCompleted)
		if !ok {
			log.Printf("unexpected event type for ScrapingCompleted handler: %T", payload)
			return
		}
		statistics := &statisticsEntity.ScrapingStatistics{
			Found:   e.Found,
			Scraped: e.Scraped,
			New:     e.New,
			OnSale:  e.OnSale,
		}
		if e.Sample != nil {
			statistics.Sample = &statisticsEntity.SampleProduct{
				Name:        e.Sample.Name,
				Price:       e.Sample.Price,
				Category:    e.Sample.Category,
				Marketplace: e.Sample.Marketplace,
			}
		}
		statisticsRepository.Publish(statistics)
	})

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
