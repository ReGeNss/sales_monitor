package scraper

import (
	"log"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/service"

	"github.com/playwright-community/playwright-go"
)

type ScraperService interface {
	Scrape() ([]*entity.ScrapedProducts, error)
}

type scraperServiceImpl struct {
	configuration  []ScraperConfig
	productService service.ProductService
}

func NewScraperService(configuration []ScraperConfig, productService service.ProductService) ScraperService {
	return &scraperServiceImpl{
		configuration:  configuration,
		productService: productService,
	}
}

func (s *scraperServiceImpl) Scrape() ([]*entity.ScrapedProducts, error) {
	scrapedProducts := []*entity.ScrapedProducts{}
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
			scrapedProducts = append(scrapedProducts, &entity.ScrapedProducts{
				Products:        products,
				MarketplaceName: scraperConfig.MarketplaceName,
				Category:        scrapingContent.Category,
			})
			log.Printf("found %d products", len(products))
		}
	}
	pw.Stop()
	return scrapedProducts, nil
}
