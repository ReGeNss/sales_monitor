package entity

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Product struct {
	ID              string
	Name            string
	RegularPrice    float64
	DiscountedPrice float64
	Image           string
}

func NewProduct(
	name string,
	regularPrice string,
	discountedPrice string,
	image string,
) *Product {
	return &Product{
		Name:            name,
		RegularPrice:    parsePrice(regularPrice),
		DiscountedPrice: parsePrice(discountedPrice),
		Image:           image,
	}
}

func parsePrice(priceText string) float64 {
	re := regexp.MustCompile(`[^\d.,]`)
	cleaned := re.ReplaceAllString(priceText, "")

	cleaned = strings.Replace(cleaned, ",", ".", -1)

	price, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		log.Printf("could not parse price '%s': %v", priceText, err)
		return 0.0
	}

	return price
}
