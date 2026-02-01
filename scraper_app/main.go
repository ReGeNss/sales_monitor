package main

import (
	"flag"
	"log"
	"strings"
	config "sales_monitor/internal/scheduler_config"
	scraper "sales_monitor/scraper_app/feature/scraper/domain/entity"
	scrapers "sales_monitor/scraper_app/feature/scraper/service/scrapers"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"github.com/joho/godotenv"
)

func main() {
	configPathFlag := flag.String("config", "", "Path to scraper config YAML")
	jobIDFlag := flag.String("job-id", "", "Job ID to execute")
	flag.Parse()

	loadEnv()

	jobID := strings.TrimSpace(*jobIDFlag)

	if jobID == "" {
		log.Fatalf("Job ID is required")
	}

	if err := runConfigJob(strings.TrimSpace(*configPathFlag), jobID); err != nil {
		log.Fatalf("Error scraping products: %v", err)
	}
}
	
func loadEnv() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalln("Error loading .env file")
	}
}

func runConfigJob(configPath string, jobID string) error {
	cfg := config.LoadConfig(configPath)

	job, err := cfg.FindJob(jobID)
	if err != nil {
		return err
	}

	scraperInstance, err := scrapers.GetScraperByShopName(job.ShopID)
	if err != nil {
		return err
	}

	plan := buildPlan(job, scraperInstance)
	return Run(plan)
}

func buildPlan(job *config.ResolvedJob, scraperInstance scraper.Scraper) scraper.ScrapingPlan {
	var differentiation *entity.ProductDifferentiationEntity
	if len(job.Differentiation) > 0 {
		differentiation = &entity.ProductDifferentiationEntity{
			Elements: job.Differentiation,
		}
	}

	return scraper.ScrapingPlan{
		Categories: []scraper.ScrapingCategory{
			{
				Category:                     job.CategoryName,
				WordsToIgnore:                job.WordsToIgnore,
				ProductDifferentiationEntity: differentiation,
				ScrapersConfigs: []scraper.ScraperConfig{
					{
						URLs:    job.URLs,
						Scraper: scraperInstance,
					},
				},
			},
		},
	}
}
