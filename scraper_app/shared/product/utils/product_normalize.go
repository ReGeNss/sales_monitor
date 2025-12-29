package utils

import (
	"regexp"
	"slices"
	"strings"
)

func NormalizeProductName(name string) string {
	loweredName := strings.ToLower(name)

	specialCharactersRegex := regexp.MustCompile(`[^\p{L}\p{N}\s]`)
	cleanedName := specialCharactersRegex.ReplaceAllString(loweredName, "")
	
	words := strings.Fields(cleanedName)

	slices.SortFunc(words, func(a, b string) int {
		return strings.Compare(a, b)
	})

	sortedName := strings.Join(words, " ")

	gramsRegex := regexp.MustCompile(`(\d+)\s*(грам|гр|г)\.*`)
	gramsFormatted := gramsRegex.ReplaceAllString(sortedName, "${1}гр")
	
	kilogramRegex := regexp.MustCompile(`(\d+)\s*(кг|кілограм|кіло|кг|кіло)\.*`)
	kilogramFormatted := kilogramRegex.ReplaceAllString(gramsFormatted, "${1}кг")

	return kilogramFormatted
}