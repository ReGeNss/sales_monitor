package dto

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	regexps "sales_monitor/scraper_app/core/regexp"
)

type ScrapedProductDto struct {
	Name         string
	RegularPrice float64
	SpecialPrice float64
	ImageURL     string
	URL          string
	BrandName    string
	Volume       string
	Weight       string
}

func CreateScrapedProductDto(
	name string,
	regularPrice string,
	specialPrice string,
	imageURL string,
	url string,
) *ScrapedProductDto {
	return &ScrapedProductDto{
		Name:         name,
		RegularPrice: parsePrice(regularPrice),
		SpecialPrice: parsePrice(specialPrice),
		ImageURL:     imageURL,
		URL:          url,
	}
}

func parsePrice(priceText string) float64 {	
	re := regexp.MustCompile(`[^\d.,]`)
	cleaned := re.ReplaceAllString(priceText, "")

	cleaned = strings.Replace(cleaned, ",", ".", -1)

	price, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		log.Printf("could not parse price '%s': %v", priceText, err)
		return 0
	}

	return price
}

func (s *ScrapedProductDto) ScraperSetVolumeOrWeight(name string) error {
	formattedValue, isVolume, err := getVolumeOrWeightFromName(name)

	if err != nil {
		return err
	}
	if isVolume {
		s.Volume = strconv.FormatFloat(formattedValue, 'f', 3, 64)
	} else {
		s.Weight = strconv.FormatFloat(formattedValue, 'f', 3, 64)
	}
	return nil
}

func getVolumeOrWeightFromName(name string) (float64, bool, error) {
	name = strings.ReplaceAll(name, ",", ".")
	cleanDecimalRegex := regexp.MustCompile(regexps.WithoutDecimalRegex)
	gramsRegex := regexp.MustCompile(regexps.GramsRegex)

	grams := gramsRegex.FindString(name)
	if grams!= "" {
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
	if  volume != "" {
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
	if  volume != "" {
		cleaned := strings.Join(cleanDecimalRegex.FindAllString(volume, -1), "")
		value, err := strconv.ParseFloat(cleaned, 64)
		if err != nil {
			return 0, false, err
		}
		return value, true, nil
	}

	return 0, false, nil
}
