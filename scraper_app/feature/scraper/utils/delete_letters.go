package utils

import (
	"log"
	"regexp"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"strconv"
	"strings"
)

func ScraperSetVolumeOrWeight(text string, product *entity.ScrapedProduct) {
	loweredText := strings.ToLower(text)
	isVolume := strings.Contains(loweredText, "л") 
	shouldDivideBy1000 := !(strings.Contains(loweredText, "кг") || strings.Contains(loweredText, "л"))

	re := regexp.MustCompile(`[^\d.,]`)
	cleaned := re.ReplaceAllString(text, "")

	value, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		log.Printf("could not parse value: %v", err)
		return
	}

	if shouldDivideBy1000 {
		value /= 1000
	}
	
	formattedValue := strconv.FormatFloat(value, 'f', 3, 64)
	if isVolume {
		product.Volume = formattedValue
	} else {
		product.Weight = formattedValue
	}
}