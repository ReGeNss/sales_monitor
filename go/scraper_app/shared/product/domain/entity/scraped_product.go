package entity

import (
	"fmt"
	"regexp"
	regexps "sales_monitor/scraper_app/core/regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type ScrapedProducts struct {
	Products        []*ScrapedProduct
	MarketplaceName string
}

type ScrapedProduct struct {
	Name         string
	RegularPrice float64
	SpecialPrice float64
	Image        string
	BrandName    string
	URL          string
	Volume       string
	Weight       string
}

func NewScrapedProduct(
	name string,
	regularPrice float64,
	specialPrice float64,
	image string,
	url string,
	brandName string,
	volume string,
	weight string,
	wordsToIgnore []string,
) (*ScrapedProduct, error) {
	var validName = strings.TrimSpace(name)
	if validName == "" {
		return nil, nil // TODO: create domain error
	}

	validName = replaceIgnoredWords(validName, wordsToIgnore)
	
	var validBrandName = strings.TrimSpace(brandName)
	if validBrandName == "" {
		return nil, nil // TODO: create domain error
	}

	var validUrl = strings.TrimSpace(url)
	if validUrl == "" {
		return nil, nil // TODO: create domain error
	}

	return &ScrapedProduct{
		Name:         validName,
		RegularPrice: regularPrice,
		SpecialPrice: specialPrice,
		Image:        image,
		BrandName:    validBrandName,
		URL:          validUrl,
		Volume:       volume,
		Weight:       weight,
	}, nil
}

func (s *ScrapedProduct) GetFingerprint(wordsToIgnore []string) *string {
	wordsToIgnore = append(wordsToIgnore, s.BrandName)
	loweredName := strings.ToLower(s.Name)

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

func replaceIgnoredWords(title string, wordsToIgnore []string) string {
    pattern := fmt.Sprintf("(?i)(%s)", strings.Join(wordsToIgnore, "|")) 
    re := regexp.MustCompile(pattern)
    
    result := re.ReplaceAllString(title, "")
    
    return result
}