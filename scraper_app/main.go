package main

import (
	"log"
	"sales_monitor/internal/db"
	// scraper "sales_monitor/scraper_app/feature/scraper/service"
	// "sales_monitor/scraper_app/feature/scraper/service/scrapers/fora"
	// "sales_monitor/scraper_app/feature/scraper/service/scrapers/silpo"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/repository"
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
	// 		// 	URLs:    []string{"https://www.atbmarket.com/promo/economy"},
	// 		// 	Scraper: atb.AtbScraper,
	// 		// },
	// 		{
	// 			URLs:    []string{"https://fora.ua/category/chypsy-2735"},
	// 			Scraper: fora.ForaScraper,
	// 		},
	// 		{
	// 			URLs:    []string{"https://silpo.ua/category/kartopliani-chypsy-5021"},
	// 			Scraper: silpo.SilpoScraper,
	// 		},
	// 	},
	// 	nil,
	// )

	// scrapedProducts, err := scraperService.Scrape()
	// if err != nil {
	// 	log.Fatalf("Error scraping products: %v", err)
	// }

	productService := service.NewProductService(repository.NewProductRepository(db.GetDB()))
	productService.ProcessProducts([]*entity.ScrapedProduct{
		{
			Name: "Чипси Lay's картопляні зі смаком макарони з сиром",
			RegularPrice: 39.99,
			DiscountedPrice: 65.99,
			Image: "https://images.silpo.ua/products/300x300/webp/09db9b6c-9b5b-4aca-8f87-c3340728f533.png",
		},
	})
}
