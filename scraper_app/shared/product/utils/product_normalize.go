package utils

import (
	"regexp"
	"slices"
	"strings"
)

func NormalizeProductName(name string, brand string, category string) string {
	loweredName := strings.ToLower(name)
	loweredBrand := strings.ToLower(brand)
	loweredCategory := strings.ToLower(category)

	gramsRegex := regexp.MustCompile(`(\d+)\s*(грам|гр|г)\.*`)
	gramsFormatted := gramsRegex.ReplaceAllString(loweredName, "${1}гр")
	
	kilogramRegex := regexp.MustCompile(`(\d+)\s*(кг|кілограм|кіло|кг|кіло)\.*`)
	kilogramFormatted := kilogramRegex.ReplaceAllString(gramsFormatted, "${1}кг")

	removedBrand := strings.ReplaceAll(kilogramFormatted, loweredBrand, "")
	removedCategory := strings.ReplaceAll(removedBrand, loweredCategory, "")

	specialCharactersRegex := regexp.MustCompile(`[^\p{L}\p{N}\s]`)
	cleanedName := specialCharactersRegex.ReplaceAllString(removedCategory, "")
	
	words := strings.Fields(cleanedName)

	slices.SortFunc(words, func(a, b string) int {
		return strings.Compare(a, b)
	})

	normalizedName := strings.Join(words, " ")


	return normalizedName
}