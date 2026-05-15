package event

type ScrapingCompleted struct {
	Found   int
	Scraped int
	New     int
	OnSale  int
	Sample  *SampleProduct
}

type SampleProduct struct {
	Name        string
	Price       float64
	Category    string
	Marketplace string
}
