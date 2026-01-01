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
	shouldMultiplyBy1000 := !(strings.Contains(loweredText, "кг") || strings.Contains(loweredText, "мл"))

	re := regexp.MustCompile(`[^\d.,]`)
	cleaned := re.ReplaceAllString(text, "")

	value, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		log.Printf("could not parse value: %v", err)
		return
	}

	if shouldMultiplyBy1000 {
		value *= 1000
	}
	
	formattedValue := strconv.Itoa(int(value))
	if isVolume {
		product.Volume = formattedValue
	} else {
		product.Weight = formattedValue
	}
}