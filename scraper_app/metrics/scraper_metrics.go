package metrics

import (
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const jobName = "scraper"

type ScrapingMetrics struct {
	Found   int
	Scraped int
	New     int
	OnSale  int
}

func PushToPrometheus(metrics ScrapingMetrics) {
	gatewayURL := os.Getenv("PUSHGATEWAY_URL")
	if gatewayURL == "" {
		log.Printf("PUSHGATEWAY_URL not set, skipping metrics push")
		return
	}

	pusher := push.New(gatewayURL, jobName)

	foundGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "scraper_products_found_total",
		Help: "Total number of products found in catalog",
	})
	foundGauge.Set(float64(metrics.Found))
	pusher.Collector(foundGauge)

	scrapedGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "scraper_products_scraped_total",
		Help: "Total number of products scraped with details",
	})
	scrapedGauge.Set(float64(metrics.Scraped))
	pusher.Collector(scrapedGauge)

	newGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "scraper_products_new_total",
		Help: "Total number of new products (not in cache)",
	})
	newGauge.Set(float64(metrics.New))
	pusher.Collector(newGauge)

	onSaleGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "scraper_products_on_sale_total",
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
