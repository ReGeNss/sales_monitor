package metrics

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const jobName = "scraper"
const maxLabelLength = 128

type ScrapingMetrics struct {
	Found   int
	Scraped int
	New     int
	OnSale  int

	SampleProductName   string  
	SampleProductPrice  float64 
	SampleCategory      string  
	SampleMarketplace   string  
}

func truncateLabel(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func PushToPrometheus(metrics ScrapingMetrics) {
	gatewayURL := os.Getenv("PUSHGATEWAY_URL")
	if gatewayURL == "" {
		log.Printf("PUSHGATEWAY_URL not set, skipping metrics push")
		return
	}

	pusher := push.New(gatewayURL, jobName)

	if metrics.SampleProductName != "" {
		pusher.Grouping("sample_product", truncateLabel(metrics.SampleProductName, maxLabelLength))
	}
	if metrics.SampleProductPrice > 0 {
		pusher.Grouping("sample_price", fmt.Sprintf("%.2f", metrics.SampleProductPrice))
	}
	if metrics.SampleCategory != "" {
		pusher.Grouping("category", truncateLabel(metrics.SampleCategory, maxLabelLength))
	}
	if metrics.SampleMarketplace != "" {
		pusher.Grouping("marketplace", truncateLabel(metrics.SampleMarketplace, maxLabelLength))
	}

	foundGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "scraper_products_found",
		Help: "Total number of products found in catalog",
	})
	foundGauge.Set(float64(metrics.Found))
	pusher.Collector(foundGauge)

	scrapedGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "scraper_products_scraped",
		Help: "Total number of products scraped with details",
	})
	scrapedGauge.Set(float64(metrics.Scraped))
	pusher.Collector(scrapedGauge)

	newGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "scraper_products_new",
		Help: "Total number of new products (not in cache)",
	})
	newGauge.Set(float64(metrics.New))
	pusher.Collector(newGauge)

	onSaleGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "scraper_products_on_sale",
		Help: "Total number of products on sale",
	})
	onSaleGauge.Set(float64(metrics.OnSale))
	pusher.Collector(onSaleGauge)

	if err := pusher.Push(); err != nil {
		log.Printf("failed to push metrics to Pushgateway: %v", err)
		return
	}
	log.Printf("metrics pushed to Prometheus Pushgateway successfully")
}
