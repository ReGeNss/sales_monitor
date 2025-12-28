package scraper

import (
	"log"
	"sales_monitor/scraper_app/shared/product/service"
	"github.com/playwright-community/playwright-go"
)

type ScraperService interface {
	Scrape()
}

type scraperServiceImpl struct {
	scrapers       []ScraperConfig
	productService service.ProductService
}

func NewScraperService(scrapers []ScraperConfig, productService service.ProductService) ScraperService {
	return &scraperServiceImpl{
		scrapers:       scrapers,
		productService: productService,
	}
}

func (s *scraperServiceImpl) Scrape() {

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

	for _, scraperConfig := range s.scrapers {
		for _, url := range scraperConfig.URLs {
			log.Printf("scraping %s", url)
			products := scraperConfig.Scraper(browser, url)
			log.Printf("found %d products", len(products))
		}
	}
	pw.Stop()
}
