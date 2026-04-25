package service

import (
	"fmt"
	"log"
	config "sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/domain/gateway"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/service"
)

type ScraperService interface {
	Scrape() (map[string]*config.ScrapingResult, error)
}

type scraperServiceImpl struct {
	configuration               config.ScrapingPlan
	productService              service.ProductService
	cachedScrapedProductService CachedScrapedProductService
	scraperFactory              gateway.ScraperFactory
	resultStorage               gateway.ResultStorage
	metricsPublisher            gateway.MetricsPublisher
	errorLogger                 gateway.ErrorLogger
}

func NewScraperService(
	configuration config.ScrapingPlan,
	productService service.ProductService,
	cachedScrapedProductService CachedScrapedProductService,
	scraperFactory gateway.ScraperFactory,
	resultStorage gateway.ResultStorage,
	metricsPublisher gateway.MetricsPublisher,
	errorLogger gateway.ErrorLogger,
) ScraperService {
	return &scraperServiceImpl{
		configuration:               configuration,
		productService:              productService,
		cachedScrapedProductService: cachedScrapedProductService,
		scraperFactory:              scraperFactory,
		resultStorage:               resultStorage,
		metricsPublisher:            metricsPublisher,
		errorLogger:                 errorLogger,
	}
}

type scrapeTotals struct {
	found, scraped, new, onSale int
}

func (t *scrapeTotals) add(other scrapeTotals) {
	t.found += other.found
	t.scraped += other.scraped
	t.new += other.new
	t.onSale += other.onSale
}

func (s *scraperServiceImpl) Scrape() (map[string]*config.ScrapingResult, error) {
	scrapedProducts := map[string]*config.ScrapingResult{}
	var totals scrapeTotals

	for _, category := range s.configuration.Categories {
		categoryTotals, err := s.scrapeCategory(category, scrapedProducts)
		if err != nil {
			return nil, err
		}
		totals.add(categoryTotals)
	}

	s.resultStorage.Save(scrapedProducts, s.categoryNames())

	s.metricsPublisher.Publish(gateway.ScrapingMetrics{
		Found:   totals.found,
		Scraped: totals.scraped,
		New:     totals.new,
		OnSale:  totals.onSale,
	}, scrapedProducts)

	return scrapedProducts, nil
}

func (s *scraperServiceImpl) scrapeCategory(
	category config.ScrapingCategory,
	out map[string]*config.ScrapingResult,
) (scrapeTotals, error) {
	var totals scrapeTotals

	for _, scraperConfig := range category.ScrapersConfigs {
		for _, url := range scraperConfig.URLs {
			scraper, err := s.scraperFactory.Get(scraperConfig.ShopID)
			if err != nil {
				return totals, fmt.Errorf("get scraper for shop %s: %w", scraperConfig.ShopID, err)
			}

			group, urlTotals := s.scrapeURL(scraper, url, category)
			appendScrapedProducts(out, category, group)
			totals.add(urlTotals)
		}
	}
	return totals, nil
}

func (s *scraperServiceImpl) scrapeURL(
	scraper gateway.Scraper,
	url string,
	category config.ScrapingCategory,
) (*entity.ScrapedProducts, scrapeTotals) {
	marketplaceName := scraper.GetMarketplaceName()

	cachedProducts, err := s.cachedScrapedProductService.GetCachedScrapedProducts(marketplaceName, category.Category)
	if err != nil {
		log.Printf("error getting cached scraped products: %v", err)
		cachedProducts = nil
	}

	result := scraper.Scrape(url, cachedProducts, category.WordsToIgnore)

	group := &entity.ScrapedProducts{
		Products:        result.Products,
		MarketplaceName: marketplaceName,
	}

	totals := scrapeTotals{
		found:   result.FoundCount,
		scraped: len(result.Products),
		new:     result.NewCount,
		onSale:  countProductsOnSale(result.Products),
	}
	return group, totals
}

func appendScrapedProducts(
	out map[string]*config.ScrapingResult,
	category config.ScrapingCategory,
	group *entity.ScrapedProducts,
) {
	existing, ok := out[category.Category]
	if !ok {
		out[category.Category] = &config.ScrapingResult{
			ScrapedProducts:              []*entity.ScrapedProducts{group},
			ProductDifferentiationEntity: category.ProductDifferentiationEntity,
		}
		return
	}
	existing.ScrapedProducts = append(existing.ScrapedProducts, group)
}

func (s *scraperServiceImpl) categoryNames() []string {
	names := make([]string, 0, len(s.configuration.Categories))
	for _, c := range s.configuration.Categories {
		names = append(names, c.Category)
	}
	return names
}

func countProductsOnSale(products []*entity.ScrapedProduct) int {
	count := 0
	for _, p := range products {
		if p.SpecialPrice > 0 && p.RegularPrice < p.SpecialPrice {
			count++
		}
	}
	return count
}
