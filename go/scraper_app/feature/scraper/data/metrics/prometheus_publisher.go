package metrics

import (
	"sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/domain/gateway"
)

type prometheusPublisher struct{}

func NewPrometheusPublisher() gateway.MetricsPublisher {
	return &prometheusPublisher{}
}

func (p *prometheusPublisher) Publish(m gateway.ScrapingMetrics, results map[string]*entity.ScrapingResult) {
	out := ScrapingMetrics{
		Found:   m.Found,
		Scraped: m.Scraped,
		New:     m.New,
		OnSale:  m.OnSale,
	}

	if name, price, category, marketplace, ok := extractSample(results); ok {
		out.SampleProductName = name
		out.SampleProductPrice = price
		out.SampleCategory = category
		out.SampleMarketplace = marketplace
	}

	PushToPrometheus(out)
}

func extractSample(results map[string]*entity.ScrapingResult) (name string, price float64, category, marketplace string, ok bool) {
	for cat, result := range results {
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
			price = p.RegularPrice()
			if p.SpecialPrice() > 0 {
				price = p.SpecialPrice()
			}
			return p.Name(), price, cat, sp.MarketplaceName, true
		}
	}
	return "", 0, "", "", false
}
