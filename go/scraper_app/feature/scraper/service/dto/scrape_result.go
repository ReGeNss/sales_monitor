package dto

type ScrapeResult struct {
	Products []*ScrapedProductDto
	FoundCount int
	NewCount int
}
