package service

import (
	"log"
	config "sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/domain/gateway"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/service"

	"github.com/playwright-community/playwright-go"
)

type ScraperService interface {
	Scrape() (map[string]*config.ScrapingResult, error)
}

type scraperServiceImpl struct {
	configuration               config.ScrapingPlan
	productService              service.ProductService
	cachedScrapedProductService CachedScrapedProductService
	resultStorage               gateway.ResultStorage
	metricsPublisher            gateway.MetricsPublisher
}

func NewScraperService(
	configuration config.ScrapingPlan,
	productService service.ProductService,
	cachedScrapedProductService CachedScrapedProductService,
	resultStorage gateway.ResultStorage,
	metricsPublisher gateway.MetricsPublisher,
) ScraperService {
	return &scraperServiceImpl{
		configuration:               configuration,
		productService:              productService,
		cachedScrapedProductService: cachedScrapedProductService,
		resultStorage:               resultStorage,
		metricsPublisher:            metricsPublisher,
	}
}

func (s *scraperServiceImpl) Scrape() (map[string]*config.ScrapingResult, error) {
	scrapedProducts := map[string]*config.ScrapingResult{}
	var totalFound, totalScraped, totalNew, totalOnSale int

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(
		playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(true),
			Args: []string{
				"--disable-blink-features=AutomationControlled",
				"--disable-dev-shm-usage",
			},
		},
	)

	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	for _, scrapingCategory := range s.configuration.Categories {
		for _, scraperConfig := range scrapingCategory.ScrapersConfigs {
			for _, url := range scraperConfig.URLs {
				scraper, err := scrapers.GetScraperByShopName(scraperConfig.ShopID, browser)
				if err != nil {
					log.Fatalf("could not get scraper for shop %s: %v", scraperConfig.ShopID, err)
				}

				cachedProducts, err := s.cachedScrapedProductService.GetCachedScrapedProducts(scraper.GetMarketplaceName(), scrapingCategory.Category)
				if err != nil {
					log.Printf("error getting cached scraped products: %v", err)
					cachedProducts = nil
				}

				result := scraper.Scrape(
					browser,
					url,
					cachedProducts,
				)

				products := []*entity.ScrapedProduct{}

				for _, p := range result.Products {
					product, _ := entity.NewScrapedProduct(
						p.Name,
						p.RegularPrice,
						p.SpecialPrice,
						p.ImageURL,
						p.URL,
						p.BrandName,
						p.Volume,
						p.Weight,
						scrapingCategory.WordsToIgnore,
					)
					if product != nil {
						products = append(products, product)
					}
				}

				if scrapedProducts[scrapingCategory.Category] == nil {
					scrapedProducts[scrapingCategory.Category] = &config.ScrapingResult{
						ScrapedProducts: []*entity.ScrapedProducts{{
							Products:        products,
							MarketplaceName: scraper.GetMarketplaceName(),
						}},
						ProductDifferentiationEntity: scrapingCategory.ProductDifferentiationEntity,
					}
				} else {
					scrapedProducts[scrapingCategory.Category].ScrapedProducts = append(scrapedProducts[scrapingCategory.Category].ScrapedProducts, &entity.ScrapedProducts{
						Products:        products,
						MarketplaceName: scraper.GetMarketplaceName(),
					})
				}

				validProducts := []*entity.ScrapedProduct{}
				for _, p := range products {
					if p != nil {
						validProducts = append(validProducts, p)
					}
				}

				onSale := countProductsOnSale(validProducts)
				totalFound += result.FoundCount
				totalScraped += len(validProducts)
				totalNew += result.NewCount
				totalOnSale += onSale
			}
		}
	}

	pw.Stop()

	scrapedCategories := []string{}
	for _, category := range s.configuration.Categories {
		scrapedCategories = append(scrapedCategories, category.Category)
	}

	s.resultStorage.Save(scrapedProducts, scrapedCategories)

	s.metricsPublisher.Publish(gateway.ScrapingMetrics{
		Found:   totalFound,
		Scraped: totalScraped,
		New:     totalNew,
		OnSale:  totalOnSale,
	}, scrapedProducts)

	return scrapedProducts, nil
}

func countProductsOnSale(products []*entity.ScrapedProduct) int {
	count := 0
	for _, p := range products {
		if p.SpecialPrice > 0 && p.RegularPrice < p.SpecialPrice {
			count++
		}
	}
	return count
}
