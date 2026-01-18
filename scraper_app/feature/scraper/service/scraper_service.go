package service

import (
	"fmt"
	"log"
	"os"
	config "sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/service"
	"sales_monitor/scraper_app/utils"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

type ScraperService interface {
	Scrape() (map[string]*config.ScrapingResult, error)
}

type scraperServiceImpl struct {
	configuration  config.ScrapingPlan
	productService service.ProductService
	cachedScrapedProductService CachedScrapedProductService
}

func NewScraperService(
	configuration config.ScrapingPlan, 
	productService service.ProductService, 
	cachedScrapedProductService CachedScrapedProductService,
	) ScraperService {
	return &scraperServiceImpl{
		configuration:  configuration,
		productService: productService,
		cachedScrapedProductService: cachedScrapedProductService,
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
				cachedProducts, err := s.cachedScrapedProductService.GetCachedScrapedProducts(scraperConfig.Scraper.GetMarketplaceName(), scrapingCategory.Category)
				if err != nil {
					log.Printf("error getting cached scraped products: %v", err)
					cachedProducts = nil
				}
				
				products := scraperConfig.Scraper.Scrape(
					browser, 
					url, 
					scrapingCategory.WordsToIgnore, 
					cachedProducts,
				)

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

	scrapedCategories := []string{}
	for _, category := range s.configuration.Categories {
		scrapedCategories = append(scrapedCategories, category.Category)
	}

	utils.SaveToJsonFile(
		&scrapedProducts, 
		fmt.Sprintf("%s/scraped_%s_%s.json", 
			os.Getenv("SCRAPED_DATA_FOLDER"),
			strings.Join(scrapedCategories, " "),
			time.Now().Format(time.DateTime)),
	)
	fmt.Printf("scraped %d products GOOOL ___________\n", len(scrapedProducts))
	return scrapedProducts, nil
}
