package service

import (
	"fmt"
	"regexp"
	"sales_monitor/scraper_app/feature/product/domain/entity"

	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
)

const similarityMatchThreshold = 91

type ProductMatcher interface {
	PickBestMatch(
		fingerprint string,
		candidates []*entity.Product,
		differentiation *entity.ProductDifferentiationEntity,
	) (*entity.Product, bool)
}

type productMatcher struct{}

func NewProductMatcher() ProductMatcher {
	return &productMatcher{}
}

func (m *productMatcher) PickBestMatch(
	fingerprint string,
	candidates []*entity.Product,
	differentiation *entity.ProductDifferentiationEntity,
) (*entity.Product, bool) {
	bestSimilarity := 0
	var best *entity.Product

	for _, c := range candidates {
		if c == nil || c.NameFingerprint == nil {
			continue
		}
		similarity := fuzzy.TokenSortRatio(fingerprint, *c.NameFingerprint)
		if similarity > bestSimilarity {
			bestSimilarity = similarity
			best = c
		}
	}

	if best == nil || bestSimilarity < similarityMatchThreshold {
		return nil, false
	}

	if !differentiates(fingerprint, *best.NameFingerprint, differentiation) {
		return nil, false
	}

	return best, true
}

func differentiates(fingerprint1, fingerprint2 string, d *entity.ProductDifferentiationEntity) bool {
	if d == nil {
		return true
	}
	for _, element := range d.Elements {
		match1 := matchesAny(fingerprint1, element)
		match2 := matchesAny(fingerprint2, element)
		if match1 == match2 {
			continue
		}
		return false
	}
	return true
}

func matchesAny(fingerprint string, elements []string) bool {
	for _, e := range elements {
		escaped := regexp.QuoteMeta(e)
		re := regexp.MustCompile(fmt.Sprintf(`\s*%[1]s\*|\s*%[1]s\s*`, escaped))
		if re.MatchString(fingerprint) {
			return true
		}
	}
	return false
}
