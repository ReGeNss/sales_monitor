package model

type ScrapedProduct struct {
	Name         string
	RegularPrice float64
	SpecialPrice float64
	ImageURL     string
	URL          string
	BrandName    string
	Volume       string
	Weight       string
}

type ScrapeResult struct {
	Products   []*ScrapedProduct
	FoundCount int
	NewCount   int
}
