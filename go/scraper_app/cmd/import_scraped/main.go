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
	"sales_monitor/scraper_app/shared/product/service/usecase"
	"sales_monitor/scraper_app/utils"
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
	productRepo := repository.NewProductRepository(gormDB)
	categoryRepo := repository.NewCategoryRepository(gormDB)
	brandRepo := repository.NewBrandRepository(gormDB)
	marketplaceRepo := repository.NewMarketplaceRepository(gormDB)
	priceRepo := repository.NewPriceRepository(gormDB)
	matcher := domainservice.NewProductMatcher()

	productService := service.NewProductService(
		usecase.NewResolveCategoryUseCase(categoryRepo),
		usecase.NewResolveMarketplaceUseCase(marketplaceRepo),
		usecase.NewAssignBrandsUseCase(brandRepo),
		usecase.NewResolveBrandUseCase(brandRepo),
		usecase.NewResolveProductUseCase(productRepo, matcher),
		usecase.NewRecordPriceUseCase(marketplaceRepo, priceRepo),
		marketplaceRepo,
		utils.NewEventBus(),
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
