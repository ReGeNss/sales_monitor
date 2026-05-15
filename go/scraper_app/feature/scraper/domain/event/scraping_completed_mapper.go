package event

import "sales_monitor/scraper_app/feature/scraper/domain/entity"

func NewSampleProduct(results map[string]*entity.ScrapingResult) *SampleProduct {
	for category, result := range results {
		if result == nil {
			continue
		}
		for _, sp := range result.ScrapedProducts {
			if sp == nil || len(sp.Products) == 0 {
				continue
			}
			p := sp.Products[0]
			if p == nil || p.Name() == "" {
				continue
			}
			price := p.RegularPrice()
			if p.SpecialPrice() > 0 {
				price = p.SpecialPrice()
			}
			return &SampleProduct{
				Name:        p.Name(),
				Price:       price,
				Category:    category,
				Marketplace: sp.MarketplaceName,
			}
		}
	}
	return nil
}
