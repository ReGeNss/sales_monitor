package service

import (
	"github.com/playwright-community/playwright-go"
	"log"
	config "sales_monitor/scraper_app/feature/scraper/entity"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/service"
)

type ScraperService interface {
	Scrape() (map[string]*config.ScrapingResult, error)
}

type scraperServiceImpl struct {
	configuration  []config.ScraperConfig
	productService service.ProductService
}

func NewScraperService(configuration []config.ScraperConfig, productService service.ProductService) ScraperService {
	return &scraperServiceImpl{
		configuration:  configuration,
		productService: productService,
	}
}

func (s *scraperServiceImpl) Scrape() (map[string]*config.ScrapingResult, error) {
	scrapedProducts := map[string]*config.ScrapingResult{}
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(
		playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(false),
			Args: []string{
				"--disable-blink-features=AutomationControlled",
				"--disable-dev-shm-usage",
			},
		},
	)
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	for _, scraperConfig := range s.configuration {
		for _, scrapingContent := range scraperConfig.ScrapingContent {
			log.Printf("scraping %s", scrapingContent.URL)
			products := scraperConfig.ScraperFunction(browser, scrapingContent.URL, scrapingContent.WordsToIgnore)
			if scrapedProducts[scrapingContent.Category] == nil {
				scrapedProducts[scrapingContent.Category] = &config.ScrapingResult{
					ScrapedProducts: []*entity.ScrapedProducts{{
						Products:        products,
						MarketplaceName: scraperConfig.MarketplaceName,
					}},
					ProductDifferentiationEntity: scrapingContent.ProductDifferentiationEntity}
			} else {
				scrapedProducts[scrapingContent.Category].ScrapedProducts = append(scrapedProducts[scrapingContent.Category].ScrapedProducts, &entity.ScrapedProducts{
					Products:        products,
					MarketplaceName: scraperConfig.MarketplaceName,
				})
			}
			log.Printf("found %d products", len(products))
		}
	}
	pw.Stop()
	return scrapedProducts, nil
}
