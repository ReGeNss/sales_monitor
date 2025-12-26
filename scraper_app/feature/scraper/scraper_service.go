package scraper

import "sales_monitor/scraper/shared/product/service"

type ScraperService struct {
	scrapers []Scraper
	productService service.ProductService 
}
