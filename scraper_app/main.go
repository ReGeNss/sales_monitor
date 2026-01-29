package main

import (
	// "encoding/json"
	// "os"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"sales_monitor/internal/schedulerconfig"
	scraper "sales_monitor/scraper_app/feature/scraper/domain/entity"

	scrapers "sales_monitor/scraper_app/feature/scraper/service/scrapers"
	// "sales_monitor/scraper_app/feature/scraper/service/scrapers/atb"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/fora"
	// "sales_monitor/scraper_app/feature/scraper/service/scrapers/silpo"
	"sales_monitor/scraper_app/shared/product/domain/entity"

	"github.com/joho/godotenv"
)

func main() {
	configPathFlag := flag.String("config", "", "Path to scraper config YAML")
	jobIDFlag := flag.String("job-id", "", "Job ID to execute")
	flag.Parse()

	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	jobID := strings.TrimSpace(*jobIDFlag)
	if jobID == "" {
		jobID = strings.TrimSpace(os.Getenv("SCRAPER_JOB_ID"))
	}

	if jobID != "" {
		if err := runConfigJob(strings.TrimSpace(*configPathFlag), jobID); err != nil {
			log.Fatalf("Error scraping products: %v", err)
		}
		return
	}

	err := Run(scraper.ScrapingPlan{
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
					// {
					// URLs:    []string{"https://www.atbmarket.com/catalog/307-napoi"},
					// Scraper: &atb.AtbScraper{},
					// },
					{
						URLs:    []string{"https://fora.ua/category/solodka-voda-2483"},
						Scraper: &fora.ForaScraper{},
					},
					{
						// 	URLs: []string{
						// 		// "https://silpo.ua/category/solodka-voda-gazovana-5095/f/brand=coca-cola",
						// // 		// "https://silpo.ua/category/solodka-voda-gazovana-5095/f/brand=pepsi",
						// // 		"https://silpo.ua/category/solodka-voda-gazovana-5095/f/brand=sprite",
						// 	},
						// 	Scraper: &silpo.SilpoScraper{},
					},
				},
			},
			// {
			// 	Category:      "Cоки, нектари",
			// 	WordsToIgnore: []string{},
			// 	ProductDifferentiationEntity: &entity.ProductDifferentiationEntity{
			// 		Elements: [][]string{
			// 			{"сік"},
			// 			{"нектар"},
			// 		},
			// 	},
			// 	ScrapersConfigs: []scraper.ScraperConfig{
			// 		{
			// 			URLs:    []string{"https://www.atbmarket.com/catalog/324-soki-nektari"},
			// 			Scraper: &atb.AtbScraper{},
			// 		},
			// 		{
			// 			URLs: []string{
			// 				"https://fora.ua/category/nektary-2489",
			// 				"https://fora.ua/category/soky-2490",
			// 			},
			// 			Scraper: &fora.ForaScraper{},
			// 		},
			// 		{
			// 			URLs: []string{
			// 				"https://silpo.ua/category/soki-nektari-5096",
			// 			},
			// 			Scraper: &silpo.SilpoScraper{},
			// 		},
			// 	},
			// },
		},
	})
	if err != nil {
		log.Fatalf("Error scraping products: %v", err)
	}
}

func loadEnv() error {
	if err := godotenv.Load(); err == nil {
		return nil
	}
	if err := godotenv.Load("../.env"); err == nil {
		return nil
	}
	return godotenv.Load("../../.env")
}

func runConfigJob(configPath string, jobID string) error {
	if configPath == "" {
		configPath = strings.TrimSpace(os.Getenv("SCRAPER_CONFIG_PATH"))
	}
	if configPath == "" {
		return fmt.Errorf("SCRAPER_CONFIG_PATH is required for job execution")
	}

	cfg, err := schedulerconfig.LoadConfig(configPath)
	if err != nil {
		return err
	}

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

func buildPlan(job *schedulerconfig.ResolvedJob, scraperInstance scraper.Scraper) scraper.ScrapingPlan {
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
