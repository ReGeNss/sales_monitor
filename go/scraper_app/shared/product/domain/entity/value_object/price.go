package valueObject

import (
	"fmt"
	"log"
	"regexp"
	"sales_monitor/scraper_app/shared/product/domain/exception"
	"strconv"
	"strings"
)

type PriceValue struct {
	price float64
}

func (s *PriceValue) GetPrice() float64 {
	return s.price
}

func NewPriceValue(rawPrice string) (*PriceValue, exception.IDomainError) {
	re := regexp.MustCompile(`[^\d.,]`)
	cleaned := re.ReplaceAllString(rawPrice, "")
	cleaned = strings.Replace(cleaned, ",", ".", -1)

	price, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		log.Printf("could not parse price '%s': %v", rawPrice, err)
		return nil, exception.NewDomainError(fmt.Sprintf("could not parse price '%s': %v", rawPrice, err))
	}

	return &PriceValue{price: price}, nil;
}
