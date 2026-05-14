package entity

type ScrapingCompleted struct {
	Found   int
	Scraped int
	New     int
	OnSale  int
	Results map[string]*ScrapingResult
}
