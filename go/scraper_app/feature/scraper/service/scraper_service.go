package service

import (
	"fmt"
	"sales_monitor/scraper_app/feature/product/domain/entity"
	config "sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/domain/exception"
	"sales_monitor/scraper_app/feature/scraper/domain/gateway"
	"sales_monitor/scraper_app/feature/scraper/domain/repository"
	"sales_monitor/scraper_app/utils"
)

type ScraperService interface {
	Scrape() (map[string]*config.ScrapingResult, exception.IDomainError)
}

type scraperServiceImpl struct {
	configuration          config.ScrapingPlan
	laterScrapedRepository repository.CachedScrapedProductsRepository
	scraperFactory         gateway.ScraperFactory
	eventBus               utils.EventBus
}

func NewScraperService(
	configuration config.ScrapingPlan,
	laterScrapedRepository repository.CachedScrapedProductsRepository,
	scraperFactory gateway.ScraperFactory,
	eventBus utils.EventBus,
) ScraperService {
	return &scraperServiceImpl{
		configuration:          configuration,
		laterScrapedRepository: laterScrapedRepository,
		scraperFactory:         scraperFactory,
		eventBus:               eventBus,
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

func (s *scraperServiceImpl) Scrape() (map[string]*config.ScrapingResult, exception.IDomainError) {
	scrapedProducts := map[string]*config.ScrapingResult{}
	var totals scrapeTotals

	for _, category := range s.configuration.Categories {
		categoryTotals, err := s.scrapeCategory(category, scrapedProducts)
		if err != nil {
			return nil, exception.NewDomainError(category.Category + " scrape error")
		}
		totals.add(categoryTotals)
	}

	s.eventBus.Publish(&config.ScrapingCompleted{
		Found:   totals.found,
		Scraped: totals.scraped,
		New:     totals.new,
		OnSale:  totals.onSale,
		Results: scrapedProducts,
	})

	return scrapedProducts, nil
}

func (s *scraperServiceImpl) scrapeCategory(
	category config.ScrapingCategory,
	out map[string]*config.ScrapingResult,
) (scrapeTotals, exception.IDomainError) {
	var totals scrapeTotals

	for _, scraperConfig := range category.ScrapersConfigs {
		for _, url := range scraperConfig.URLs {
			scraper, err := s.scraperFactory.Get(scraperConfig.ShopID)
			if err != nil {
				return totals, exception.NewDomainError(fmt.Sprintf("get scraper for shop %s: %s", scraperConfig.ShopID, err))
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

	cachedProducts, err := s.laterScrapedRepository.GetCachedScrapedProducts(marketplaceName, category.Category)
	if err != nil {
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

func countProductsOnSale(products []*entity.ScrapedProduct) int {
	count := 0
	for _, p := range products {
		if p.SpecialPrice() > 0 && p.RegularPrice() < p.SpecialPrice() {
			count++
		}
	}
	return count
}
