package main

import (
	"encoding/json"
	"log"
	"os"
	"sales_monitor/internal/db"

	// scraper "sales_monitor/scraper_app/feature/scraper/service"
	// "sales_monitor/scraper_app/feature/scraper/service/scrapers/atb"
	// "sales_monitor/scraper_app/feature/scraper/service/scrapers/silpo"
	// "sales_monitor/scraper_app/feature/scraper/service/scrapers/fora"
	"sales_monitor/scraper_app/shared/product/data/repository"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/service"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db.ConnectToDB()

	// scraperService := scraper.NewScraperService(
	// 	[]scraper.ScraperConfig{
	// 		// {
	// 		// 	ScrapingContent: []scraper.ScrapingContent{
	// 		// 		{
	// 		// 			URL: "https://www.atbmarket.com/catalog/cipsi",
	// 		// 			Category: "Чипси",
	// 		// 		},
	// 		// 	},	
	// 		// 	MarketplaceName: "АТБ",
	// 		// 	Scraper: atb.AtbScraper,
	// 		// },
	// 		// {
	// 		// 	ScrapingContent: []scraper.ScrapingContent{
	// 		// 		{
	// 		// 			URL: "https://fora.ua/category/chypsy-2735",
	// 		// 			Category: "Чипси",
	// 		// 		},
	// 		// 	},
	// 		// 	Scraper: fora.ForaScraper,
	// 		// },
	// 		{
	// 			ScrapingContent: []scraper.ScrapingContent{
	// 				{
	// 					URL: "https://silpo.ua/category/kartopliani-chypsy-5021/f/brand=lay-s",
	// 					Category: "Чипси",
	// 				},
	// 			},
	// 			Scraper: silpo.SilpoScraper,
	// 			MarketplaceName: "Сільпо",
	// 		},
	// 	},
	// 	nil,
	// )

	// scrapedProducts, err := scraperService.Scrape()
	// if err != nil {
	// 	log.Fatalf("Error scraping products: %v", err)
	// }

	// INSERT_YOUR_CODE
	// // Write all scraped products to a JSON file
	// file, err := os.Create("scraped_products.json")
	// if err != nil {
	// 	log.Fatalf("Error creating JSON file: %v", err)
	// }
	// defer file.Close()

	// encoder := json.NewEncoder(file)
	// encoder.SetIndent("", "  ")
	// if err := encoder.Encode(scrapedProducts); err != nil {
	// 	log.Fatalf("Error encoding products to JSON: %v", err)
	// }


	// INSERT_YOUR_CODE

	// Read scraped_products.json
	file, err := os.Open("scraped_products.json")
	if err != nil {
		log.Fatalf("Error opening scraped_products.json: %v", err)
	}
	defer file.Close()

	var scrapedProducts []*entity.ScrapedProducts
	if err := json.NewDecoder(file).Decode(&scrapedProducts); err != nil {
		log.Fatalf("Error decoding scraped_products.json: %v", err)
	}


	productService := service.NewProductService(repository.NewProductRepository(db.GetDB()))
	productService.ProcessProducts(scrapedProducts)
}
