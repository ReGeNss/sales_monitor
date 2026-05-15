package repository

import (
	"fmt"
	"log"
	"os"
	"strings"

	"sales_monitor/scraper_app/feature/statistics/domain/entity"
	"sales_monitor/scraper_app/feature/statistics/domain/repository"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	jobName        = "scraper"
	maxLabelLength = 128
)

type prometheusPublisher struct{}

func NewPrometheusPublisher() repository.StatisticsRepository {
	return &prometheusPublisher{}
}

func (p *prometheusPublisher) Publish(statistics *entity.ScrapingStatistics) {
	if statistics == nil {
		return
	}

	gatewayURL := os.Getenv("PUSHGATEWAY_URL")
	if gatewayURL == "" {
		log.Printf("PUSHGATEWAY_URL not set, skipping statistics push")
		return
	}

	pusher := push.New(gatewayURL, jobName)

	if statistics.Sample != nil {
		if statistics.Sample.Name != "" {
			pusher.Grouping("sample_product", truncateLabel(statistics.Sample.Name, maxLabelLength))
		}
		if statistics.Sample.Price > 0 {
			pusher.Grouping("sample_price", fmt.Sprintf("%.2f", statistics.Sample.Price))
		}
		if statistics.Sample.Category != "" {
			pusher.Grouping("category", truncateLabel(statistics.Sample.Category, maxLabelLength))
		}
		if statistics.Sample.Marketplace != "" {
			pusher.Grouping("marketplace", truncateLabel(statistics.Sample.Marketplace, maxLabelLength))
		}
	}

	pusher.Collector(newGauge("scraper_products_found", "Total number of products found in catalog", statistics.Found))
	pusher.Collector(newGauge("scraper_products_scraped", "Total number of products scraped with details", statistics.Scraped))
	pusher.Collector(newGauge("scraper_products_new", "Total number of new products (not in cache)", statistics.New))
	pusher.Collector(newGauge("scraper_products_on_sale", "Total number of products on sale", statistics.OnSale))

	if err := pusher.Push(); err != nil {
		log.Printf("failed to push statistics to Pushgateway: %v", err)
		return
	}
	log.Printf("statistics pushed to Prometheus Pushgateway successfully")
}

func newGauge(name, help string, value int) prometheus.Gauge {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{Name: name, Help: help})
	gauge.Set(float64(value))
	return gauge
}

func truncateLabel(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
