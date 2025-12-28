package main

import (
	"log"
	"sales_monitor/internal/db"
	scraper "sales_monitor/scraper_app/feature/scraper/service"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/atb"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/fora"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/silpo"
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
				URLs:    []string{"https://www.atbmarket.com/promo/economy"},
				Scraper: atb.AtbScraper,
			},
			{
				URLs:    []string{"https://fora.ua/all-offers?filter_CATEGORY=(2730)"},
				Scraper: fora.ForaScraper,
			},
			{
				URLs:    []string{"https://silpo.ua/category/kartopliani-chypsy-5021"},
				Scraper: silpo.SilpoScraper,
			},
		},
		nil,
	)

	scraperService.Scrape()
}
