package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"

	"sales_monitor/internal/db"
	scraper "sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/shared/product/data/repository"
	domainservice "sales_monitor/scraper_app/shared/product/domain/service"
	"sales_monitor/scraper_app/shared/product/service"
)

func main() {
	inputPath := flag.String("input", ".logs/scraped_products.json", "Path to scraped products JSON")
	flag.Parse()

	loadEnv()
	db.ConnectToDB()

	scrapedData, err := loadScrapedData(*inputPath)
	if err != nil {
		log.Fatalf("failed to load scraped data: %v", err)
	}

	gormDB := db.GetDB()
	productService := service.NewProductService(
		repository.NewProductRepository(gormDB),
		repository.NewCategoryRepository(gormDB),
		repository.NewBrandRepository(gormDB),
		repository.NewMarketplaceRepository(gormDB),
		repository.NewPriceRepository(gormDB),
		repository.NewNotificationPublisher(db.GetRedis()),
		domainservice.NewProductMatcher(),
	)

	productService.ProcessProducts(scrapedData)
	log.Printf("processed %d categories", len(scrapedData))
}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Error loading .env file")
	}
}

func loadScrapedData(path string) (map[string]*scraper.ScrapingResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var scrapedData map[string]*scraper.ScrapingResult
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&scrapedData); err != nil {
		return nil, err
	}
	return scrapedData, nil
}
