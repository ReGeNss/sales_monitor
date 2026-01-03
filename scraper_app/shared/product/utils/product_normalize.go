package utils

import (
	"regexp"
	"slices"
	"strings"
)

func NormalizeProductName(name string, wordsToIgnore []string) string {
	loweredName := strings.ToLower(name)

	gramsRegex := regexp.MustCompile(`(\d+)\s*(грам|гр|г)\.*`)
	gramsFormatted := gramsRegex.ReplaceAllString(loweredName, "${1}гр")
	
	kilogramRegex := regexp.MustCompile(`(\d+)\s*(кг|кілограм|кіло|кг|кіло)\.*`)
	kilogramFormatted := kilogramRegex.ReplaceAllString(gramsFormatted, "${1}кг")

	cleaned := kilogramFormatted
	for _, word := range wordsToIgnore {
		cleaned = strings.ReplaceAll(cleaned, strings.ToLower(word), "")
	}

	specialCharactersRegex := regexp.MustCompile(`[^\p{L}\p{N}\s]`)
	cleanedSpecialCharacters := specialCharactersRegex.ReplaceAllString(cleaned, "")
	
	words := strings.Fields(cleanedSpecialCharacters)

	slices.SortFunc(words, func(a, b string) int {
		return strings.Compare(a, b)
	})

	normalizedName := strings.Join(words, " ")
	
	return normalizedName
}