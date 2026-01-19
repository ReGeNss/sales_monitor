package utils

import (
	"fmt"
	"log"
	"regexp"
	"sales_monitor/scraper_app/shared/product/domain/entity"
)

func ProductDifferentiator(fingerprint1 string, fingerprint2 string, productDifferentiationEntity *entity.ProductDifferentiationEntity) bool {
	if productDifferentiationEntity == nil {
		return true
	}
	for _, element := range productDifferentiationEntity.Elements {
		var match1 = false
		var match2 = false
		for _, e := range element {
			escapedE := regexp.QuoteMeta(e)
			re := regexp.MustCompile(fmt.Sprintf(`\s*%[1]s\*|\s*%[1]s\s*`, escapedE))

			if re.MatchString(fingerprint1) {
				match1 = true
				continue
			}
		}

		for _, e := range element {
			escapedE := regexp.QuoteMeta(e)
			re := regexp.MustCompile(fmt.Sprintf(`\s*%[1]s\*|\s*%[1]s\s*`, escapedE))

			if re.MatchString(fingerprint2) {
				match2 = true
				continue
			}
		}

		if match1 == match2 {
			log.Printf("match1: %v, match2: %v", fingerprint1, fingerprint2)
			continue
		}
		return false
	}

	return true
}
