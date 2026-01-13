package service

import (
	"fmt"
	"log"
	"os"
	config "sales_monitor/scraper_app/feature/scraper/entity"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/service"
	"sales_monitor/scraper_app/utils"
	"time"

	"github.com/playwright-community/playwright-go"
)

type ScraperService interface {
	Scrape() (map[string]*config.ScrapingResult, error)
}

type scraperServiceImpl struct {
	configuration  config.ScrapingPlan
	productService service.ProductService
}

func NewScraperService(configuration config.ScrapingPlan, productService service.ProductService) ScraperService {
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

	for _, scrapingCategory := range s.configuration.Categories {
		for _, scraperConfig := range scrapingCategory.ScrapersConfigs {
			for _, url := range scraperConfig.URLs {
				products := scraperConfig.Scraper.Scrape(browser, url, scrapingCategory.WordsToIgnore)
				if scrapedProducts[scrapingCategory.Category] == nil {
					scrapedProducts[scrapingCategory.Category] = &config.ScrapingResult{
						ScrapedProducts: []*entity.ScrapedProducts{{
							Products:        products,
							MarketplaceName: scraperConfig.Scraper.GetMarketplaceName(),
						}},
					}
				} else {
					scrapedProducts[scrapingCategory.Category].ScrapedProducts = append(scrapedProducts[scrapingCategory.Category].ScrapedProducts, &entity.ScrapedProducts{
						Products:        products,
						MarketplaceName: scraperConfig.Scraper.GetMarketplaceName(),
					})
				}
				log.Printf("found %d products", len(products))
			}
		}
	}

	pw.Stop()

	utils.SaveToJsonFile(&scrapedProducts, fmt.Sprintf("%s/scraped_products_%s.json", os.Getenv("SCRAPED_DATA_FOLDER"), time.Now().Format(time.DateTime)))
	return scrapedProducts, nil
}
