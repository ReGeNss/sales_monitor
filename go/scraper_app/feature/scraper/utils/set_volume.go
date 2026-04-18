package utils

import (
	"regexp"
	regexps "sales_monitor/scraper_app/core/regexp"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"strconv"
	"strings"
)

func ScraperSetVolumeOrWeight(text string, product *entity.ScrapedProduct) error {
	formattedValue, isVolume, err := GetVolumeOrWeightFromName(text)

	if err != nil {
		return err
	}
	if isVolume {
		product.Volume = strconv.FormatFloat(formattedValue, 'f', 3, 64)
	} else {
		product.Weight = strconv.FormatFloat(formattedValue, 'f', 3, 64)
	}
	return nil
}

func GetVolumeOrWeightFromName(name string) (float64, bool, error) {
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