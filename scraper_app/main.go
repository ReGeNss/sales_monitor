package main

import (
	// "encoding/json"
	"log"
	// "os"
	"sales_monitor/internal/db"
	"sales_monitor/scraper_app/core/api"
	scraper "sales_monitor/scraper_app/feature/scraper/domain/entity"
	scraper_service "sales_monitor/scraper_app/feature/scraper/service"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/atb"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/fora"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/silpo"
	"sales_monitor/scraper_app/shared/product/data/repository"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/service"
	cached_scraped_product_service "sales_monitor/scraper_app/feature/scraper/service"
	cached_scraped_product_repository "sales_monitor/scraper_app/feature/scraper/data/repository"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db.ConnectToDB()

	cachedScrapedProductService := cached_scraped_product_service.NewCachedScrapedProductService(cached_scraped_product_repository.NewCachedScrapedProductsRepository(db.GetDB()))

	scraperService := scraper_service.NewScraperService(
		scraper.ScrapingPlan{
			Categories: []scraper.ScrapingCategory{
				// {
				// 	Category:      "Чипси",
				// 	WordsToIgnore: []string{},
				// 	ProductDifferentiationEntity: &entity.ProductDifferentiationEntity{
				// 		Elements: [][]string{
				// 			{},
				// 		},
				// 	},
				// 	ScrapersConfigs: []scraper.ScraperConfig{
				// 		{
				// 			URLs:    []string{"https://www.atbmarket.com/catalog/cipsi"},
				// 			Scraper: &atb.AtbScraper{},
				// 		},
				// 		{
				// 			URLs:    []string{"https://fora.ua/category/chypsy-2735"},
				// 			Scraper: &fora.ForaScraper{},
				// 		},
				// 		{
				// 			URLs:    []string{"https://silpo.ua/category/kartopliani-chypsy-5021/f/brand=lay-s"},
				// 			Scraper: &silpo.SilpoScraper{},
				// 		},
				// 	},
				// },
				{
					Category: "Напої газовані",
					WordsToIgnore: []string{
						"безалкогольний",
						"напій",
					},
					ProductDifferentiationEntity: &entity.ProductDifferentiationEntity{
						Elements: [][]string{
							{"зб"},
							{"пет"},
						},
					},
					ScrapersConfigs: []scraper.ScraperConfig{
						{
							URLs:    []string{"https://www.atbmarket.com/catalog/307-napoi"},
							Scraper: &atb.AtbScraper{},
						},
						{
							URLs:    []string{"https://fora.ua/category/solodka-voda-2483"},
							Scraper: &fora.ForaScraper{},
						},
						{
							URLs: []string{
								"https://silpo.ua/category/solodka-voda-gazovana-5095/f/brand=coca-cola",
								"https://silpo.ua/category/solodka-voda-gazovana-5095/f/brand=pepsi",
								"https://silpo.ua/category/solodka-voda-gazovana-5095/f/brand=sprite",
							},
							Scraper: &silpo.SilpoScraper{},
						},
					},
				},
				{
					Category:      "Cоки, нектари",
					WordsToIgnore: []string{},
					ProductDifferentiationEntity: &entity.ProductDifferentiationEntity{
						Elements: [][]string{
							{"сік"},
							{"нектар"},
						},
					},
					ScrapersConfigs: []scraper.ScraperConfig{
						{
							URLs:    []string{"https://www.atbmarket.com/catalog/324-soki-nektari"},
							Scraper: &atb.AtbScraper{},
						},
						{
							URLs: []string{
								"https://fora.ua/category/nektary-2489",
								"https://fora.ua/category/soky-2490",
							},
							Scraper: &fora.ForaScraper{},
						},
						{
							URLs: []string{
								"https://silpo.ua/category/soki-nektari-5096",
							},
							Scraper: &silpo.SilpoScraper{},
						},
					},
				},
			},
		},
		nil,
		cachedScrapedProductService,
	)

	scrapedProducts, err := scraperService.Scrape()
	if err != nil {
		log.Fatalf("Error scraping products: %v", err)
	}

	// var scrapedProducts map[string]*scraper.ScrapingResult
	// file, err := os.Open("../.logs/scraped_products.json")
	// if err != nil {
	// 	log.Fatalf("Error opening scraped products file: %v", err)
	// }
	// defer file.Close()

	// err = json.NewDecoder(file).Decode(&scrapedProducts)
	// if err != nil {
	// 	log.Fatalf("Error reading scraped products from file: %v", err)
	// }

	productService := service.NewProductService(repository.NewProductRepository(db.GetDB(), api.NewHTTPClient()))
	productService.ProcessProducts(scrapedProducts)
}
