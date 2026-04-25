package entity

import (
	"fmt"
	"regexp"
	regexps "sales_monitor/scraper_app/core/regexp"
	valueObject "sales_monitor/scraper_app/shared/product/domain/entity/value_object"
	"sales_monitor/scraper_app/shared/product/domain/exception"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ScrapedProducts struct {
	Products        []*ScrapedProduct
	MarketplaceName string
}

type ScrapedProduct struct {
	name         string
	regularPrice *valueObject.PriceValue
	specialPrice *valueObject.PriceValue
	image        *valueObject.Url
	brandName    string
	url          *valueObject.Url
	volume       string
	weight       string
}

func (s *ScrapedProduct) Name() string          { return s.name }
func (s *ScrapedProduct) RegularPrice() float64 { return s.regularPrice.GetPrice() }
func (s *ScrapedProduct) SpecialPrice() float64 { return s.specialPrice.GetPrice() }
func (s *ScrapedProduct) ImageUrl() string      { return s.URL() }

func (s *ScrapedProduct) BrandName() string { return s.brandName }

func (s *ScrapedProduct) SetBrandName(brand string) exception.IDomainError {
	validBrand := strings.TrimSpace(brand)
	if validBrand == "" {
		return exception.NewDomainError("Brand name is empty")
	}
	s.brandName = validBrand
	return nil
}

func (s *ScrapedProduct) URL() string    { return s.url.Url() }
func (s *ScrapedProduct) Volume() string { return s.volume }
func (s *ScrapedProduct) Weight() string { return s.weight }

func (s *ScrapedProduct) Validate() exception.IDomainError {
	if s.name == "" {
		return exception.NewDomainError("Name is empty")
	}

	if s.url == nil {
		return exception.NewDomainError("url is empty")
	}

	if s.brandName == "" {
		return exception.NewDomainError("Brand is empty")
	}

	if s.regularPrice == nil {
		return exception.NewDomainError("Regular price is missing")
	}

	if s.specialPrice == nil {
		return exception.NewDomainError("Special price is missing")
	}

	if s.specialPrice.GetPrice() > s.regularPrice.GetPrice() {
		return exception.NewDomainError("Special price cannot be greater than regular price")
	}

	if (s.volume == "" && s.weight == "") {
		return exception.NewDomainError("One of volume or weight is required")
	}

	return nil
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
		return nil, exception.NewDomainError("Name is empty")
	}

	validName = replaceIgnoredWords(validName, wordsToIgnore)

	var validBrandName = strings.TrimSpace(brandName)
	if validBrandName == "" {
		return nil, exception.NewDomainError("Brand is empty")
	}

	validUrl, urlErr := valueObject.NewUrl(url)
	if urlErr != nil {
		return nil, urlErr
	}

	validImageUrl, urlErr := valueObject.NewUrl(image)
	if urlErr != nil {
		return nil, urlErr
	}

	validRegularPrice, err := valueObject.NewPriceValue(fmt.Sprintf("%f", regularPrice))
	if err != nil {
		return nil, err
	}

	validSpecialPrice, err := valueObject.NewPriceValue(fmt.Sprintf("%f", specialPrice))
	if err != nil {
		return nil, err
	}

	return &ScrapedProduct{
		name:         validName,
		regularPrice: validRegularPrice,
		specialPrice: validSpecialPrice,
		image:        validImageUrl,
		brandName:    validBrandName,
		url:          validUrl,
		volume:       volume,
		weight:       weight,
	}, nil
}

func CreateEmptyScrapedProduct(name string, regularPrice string, specialPrice string, imageUrl string, url string, wordsToIgnore []string) (*ScrapedProduct, exception.IDomainError) {
	validName := strings.TrimSpace(name)
	if validName == "" {
		return nil, exception.NewDomainError("Name is empty")
	}

	validName = replaceIgnoredWords(validName, wordsToIgnore)

	validRegularPrice, err := valueObject.NewPriceValue(regularPrice)
	if err != nil {
		return nil, err
	}

	validSpecialPrice, err := valueObject.NewPriceValue(specialPrice)
	if err != nil {
		return nil, err
	}

	if validSpecialPrice.GetPrice() < validRegularPrice.GetPrice() {
		return nil, exception.NewDomainError("price smaller then ")
	}

	validUrl, urlErr := valueObject.NewUrl(url)
	if urlErr != nil {
		return nil, urlErr
	}

	validImageUrl, err := valueObject.NewUrl(imageUrl)
	if err != nil {
		return nil, err
	}

	return &ScrapedProduct{
		name:         validName,
		regularPrice: validRegularPrice,
		specialPrice: validSpecialPrice,
		image:        validImageUrl,
		url:          validUrl,
	}, nil
}

func (s *ScrapedProduct) GetFingerprint(wordsToIgnore []string) *string {
	wordsToIgnore = append(wordsToIgnore, s.brandName)
	loweredName := strings.ToLower(s.name)

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

func (p *ScrapedProduct) SetVolumeOrWeight(value string) error {
	formattedValue, isVolume, err := getVolumeOrWeightFromText(value)
	if err != nil {
		return err
	}
	if isVolume {
		p.volume = strconv.FormatFloat(formattedValue, 'f', 3, 64)
	} else {
		p.weight = strconv.FormatFloat(formattedValue, 'f', 3, 64)
	}
	return nil
}

func getVolumeOrWeightFromText(name string) (float64, bool, error) {
	name = strings.ReplaceAll(name, ",", ".")
	cleanDecimalRegex := regexp.MustCompile(regexps.WithoutDecimalRegex)
	gramsRegex := regexp.MustCompile(regexps.GramsRegex)

	grams := gramsRegex.FindString(name)
	if grams != "" {
		cleaned := strings.Join(cleanDecimalRegex.FindAllString(grams, -1), "")
		value, err := strconv.ParseFloat(cleaned, 64)
		if err != nil {
			return 0, false, err
		}
		value /= 1000
		return value, false, nil
	}

	kilogramRegex := regexp.MustCompile(regexps.KilogramRegex)
	kilogram := kilogramRegex.FindString(name)
	if kilogram != "" {
		cleaned := strings.Join(cleanDecimalRegex.FindAllString(kilogram, -1), "")
		value, err := strconv.ParseFloat(cleaned, 64)
		if err != nil {
			return 0, false, err
		}
		return value, false, nil
	}

	volumeRegex := regexp.MustCompile(regexps.VolumeMilliliterRegex)
	volume := volumeRegex.FindString(name)
	if volume != "" {
		cleaned := strings.Join(cleanDecimalRegex.FindAllString(volume, -1), "")
		value, err := strconv.ParseFloat(cleaned, 64)
		if err != nil {
			return 0, false, err
		}
		value /= 1000
		return value, true, nil
	}

	volumeRegex = regexp.MustCompile(regexps.VolumeLiterRegex)
	volume = volumeRegex.FindString(name)
	if volume != "" {
		cleaned := strings.Join(cleanDecimalRegex.FindAllString(volume, -1), "")
		value, err := strconv.ParseFloat(cleaned, 64)
		if err != nil {
			return 0, false, err
		}
		return value, true, nil
	}

	return 0, false, nil
}
