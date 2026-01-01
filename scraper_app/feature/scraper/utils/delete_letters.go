package utils

import (
	"regexp"
	"strings"
)

func ScraperFormatVolumeWeight(text string) string {
	re := regexp.MustCompile(`[^\d.,]`)
	cleaned := re.ReplaceAllString(text, "")

	return strings.ReplaceAll(cleaned, ",", ".")
}