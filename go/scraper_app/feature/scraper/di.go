package scraper

import (
	"sales_monitor/scraper_app/feature/scraper/data/logging"
	"sales_monitor/scraper_app/feature/scraper/data/repository"
	scraperFactory "sales_monitor/scraper_app/feature/scraper/data/scraper"
	"sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/domain/exception"
	"sales_monitor/scraper_app/feature/scraper/domain/gateway"
	"sales_monitor/scraper_app/feature/scraper/service"
	"sales_monitor/scraper_app/utils"

	"gorm.io/gorm"
)

func NewScraperService(
	plan entity.ScrapingPlan,
	db *gorm.DB,
	eventBus utils.EventBus,
) (service.ScraperService, gateway.ScraperFactory, exception.IDomainError) {
	errorLogger := logging.NewScreenshotErrorLogger()

	factory, err := scraperFactory.NewScraperFactory(errorLogger)
	if err != nil {
		return nil, nil, err
	}

	cachedScrapedProductsRepository := repository.NewCachedScrapedProductsRepository(db)

	return service.NewScraperService(plan, cachedScrapedProductsRepository, factory, eventBus), factory, nil
}
