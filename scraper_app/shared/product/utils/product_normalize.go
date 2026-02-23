package utils

import (
	"regexp"
	regexps "sales_monitor/scraper_app/core/regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

func NormalizeProductName(name string, wordsToIgnore []string) *string {
	loweredName := strings.ToLower(name)

	gramsRegex := regexp.MustCompile(regexps.GramsRegex)
	gramsFormatted := gramsRegex.ReplaceAllString(loweredName, "")
	
	kilogramRegex := regexp.MustCompile(regexps.KilogramRegex)
	kilogramFormatted := kilogramRegex.ReplaceAllString(gramsFormatted, "")

	cleaned := kilogramFormatted
	for _, word := range wordsToIgnore {
		cleaned = strings.ReplaceAll(cleaned, strings.ToLower(word), "")
	}

	volumeRegex := regexp.MustCompile(regexps.VolumeMilliliterRegex)
	cleaned = volumeRegex.ReplaceAllString(cleaned, "")

	specialCharactersRegex := regexp.MustCompile(`[^\p{L}\p{N}\s]`)
	cleanedSpecialCharacters := specialCharactersRegex.ReplaceAllString(cleaned, "")
	
	words := strings.Fields(cleanedSpecialCharacters)

	deletedSmallWords := []string{}
	for _, word := range words {
		if utf8.RuneCountInString(word) > 1 {
			deletedSmallWords = append(deletedSmallWords, word)
		}
	}

	slices.SortFunc(deletedSmallWords, func(a, b string) int {
		return strings.Compare(a, b)
	})

	normalizedName := strings.Join(deletedSmallWords, " ")
	
	if normalizedName == "" {
		return nil
	}

	return &normalizedName
}