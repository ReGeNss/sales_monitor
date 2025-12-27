package main

import (
	scraper "sales_monitor/scraper/feature/scraper/service"
	"sales_monitor/scraper/feature/scraper/service/scrapers/fora"
)

func main() {
	scraperService := scraper.NewScraperService(
		[]scraper.ScraperConfig{
			// {
			// 	URLs: []string{"https://www.atbmarket.com/promo/economy"},
			// 	Scraper: scrapers.AtbScraper,
			// },
			{
				URLs: []string{"https://fora.ua/all-offers?filter_CATEGORY=(2730)"},
				Scraper: fora.ForaScraper,
			},
		},
		nil,
	)

	scraperService.Scrape()
}