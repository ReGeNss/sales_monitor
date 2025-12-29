package main

import (
	"log"
	"sales_monitor/internal/db"
	scraper "sales_monitor/scraper_app/feature/scraper/service"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/atb"
	// "sales_monitor/scraper_app/feature/scraper/service/scrapers/silpo"
	"sales_monitor/scraper_app/shared/product/data/repository"
	"sales_monitor/scraper_app/shared/product/service"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db.ConnectToDB()

	scraperService := scraper.NewScraperService(
		[]scraper.ScraperConfig{
			{
				ScrapingContent: []scraper.ScrapingContent{
					{
						URL: "https://www.atbmarket.com/catalog/cipsi",
						Category: "Чипси",
					},
				},
				MarketplaceName: "АТБ",
				Scraper: atb.AtbScraper,
			},
			// {
			// 	ScrapingContent: []scraper.ScrapingContent{
			// 		{
			// 			URL: "https://fora.ua/category/chypsy-2735",
			// 			Category: "Чипси",
			// 		},
			// 	},
			// 	Scraper: fora.ForaScraper,
			// },
			// {
			// 	URLs:    []string{"https://silpo.ua/category/kartopliani-chypsy-5021/f/brand=lyuks"},
			// 	Scraper: silpo.SilpoScraper,
			// },
		},
		nil,
	)

	scrapedProducts, err := scraperService.Scrape()
	if err != nil {
		log.Fatalf("Error scraping products: %v", err)
	}

	productService := service.NewProductService(repository.NewProductRepository(db.GetDB()))
	productService.ProcessProducts(scrapedProducts)
}
