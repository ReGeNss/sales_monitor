package main

import (
	"log"
	"sales_monitor/internal/db"
	"sales_monitor/scraper_app/shared/product/data/repository"
	"sales_monitor/scraper_app/shared/product/service"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/atb"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/fora"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/silpo"
	"sales_monitor/scraper_app/feature/scraper/service"
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
						URL:      "https://www.atbmarket.com/catalog/307-napoi",
						Category: "Напої газовані",
						WordsToIgnore: []string{
							"безалкогольний",
							"напій",
						},
					},
					{
						URL:      "https://www.atbmarket.com/catalog/324-soki-nektari",
						Category: "Cоки, нектари",
						WordsToIgnore: []string{},
					},
				},
				MarketplaceName: "АТБ",
				ScraperFunction: atb.AtbScraper,
			},
			{
				ScrapingContent: []scraper.ScrapingContent{
					// {
						// URL: "https://fora.ua/category/chypsy-2735",
						// Category: "Чипси",
					// },
					{
						URL: "https://fora.ua/category/solodka-voda-2483",
						Category: "Напої газовані",
						WordsToIgnore: []string{
							"безалкогольний",
							"напій",
						},
					},
					{
						URL: "https://fora.ua/category/nektary-2489",
						Category: "Cоки, нектари",
						WordsToIgnore: []string{},
					},
					{
						URL: "https://fora.ua/category/soky-2490",
						Category: "Cоки, нектари",
						WordsToIgnore: []string{},
					},
				},
				ScraperFunction: fora.ForaScraper,
				MarketplaceName: "Фора",
			},
			{
				ScrapingContent: []scraper.ScrapingContent{
					// {
					// 	URL: "https://silpo.ua/category/kartopliani-chypsy-5021/f/brand=lay-s",
					// 	Category: "Чипси",
					// 	WordsToIgnore: []string{},
					// },
					{
						URL: "https://silpo.ua/category/solodka-voda-gazovana-5095/f/brand=coca-cola",
						Category: "Напої газовані",
						WordsToIgnore: []string{
							"безалкогольний",
							"напій",
						},
					},
					{
						URL: "https://silpo.ua/category/solodka-voda-gazovana-5095/f/brand=pepsi",
						Category: "Напої газовані",
						WordsToIgnore: []string{
							"безалкогольний",
							"напій",
						},
					},
					{
						URL: "https://silpo.ua/category/solodka-voda-gazovana-5095/f/brand=sprite",
						Category: "Напої газовані",
						WordsToIgnore: []string{
							"безалкогольний",
							"напій",
						},
					},
				},
				ScraperFunction: silpo.SilpoScraper,
				MarketplaceName: "Сільпо",
			},
		},
		nil,
	)

	scrapedProducts, err := scraperService.Scrape()
	if err != nil {
		log.Fatalf("could not scrape products: %v", err)
	}

	productService := service.NewProductService(repository.NewProductRepository(db.GetDB()))
	productService.ProcessProducts(scrapedProducts)
}
