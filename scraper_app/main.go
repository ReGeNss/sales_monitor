package main

import (
	scraper "sales_monitor/scraper/feature/scraper/service"
	"sales_monitor/scraper/feature/scraper/service/scrapers"
)

func main() {
	scraperService := scraper.NewScraperService(
		[]scraper.ScraperConfig{
			{
				URLs: []string{"https://www.atbmarket.com/catalog/cipsi"},
				Scraper: scrapers.AtbScraper,
			},
		},
		nil,
	)

	scraperService.Scrape()
}