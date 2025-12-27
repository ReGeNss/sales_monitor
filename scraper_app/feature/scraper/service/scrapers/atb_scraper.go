package scrapers

import (
	"log"
	"math"
	"fmt"
	"regexp"
	"sales_monitor/scraper/shared/product/domain/entity"
	"strconv"
	"strings"
	"github.com/playwright-community/playwright-go"
)

func AtbScraper(browser playwright.Browser, url string) []entity.Product {
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	page.Goto(url)
	page.WaitForLoadState()

	countOfAllPages := getCountOfAllPages(page)

	var products []entity.Product

	for i := 1; i < countOfAllPages; i++ {
		products = append(products, getProducts(page)...)
		page.Close()
		
		context, err := browser.NewContext()
		page, err = context.NewPage()
		if err != nil {
			log.Fatalf("could not create page: %v", err)
		}
		url := fmt.Sprintf("%s?page=%d", url, i+1)
		page.Goto(url)
		page.WaitForLoadState()
	}

	products = append(products, getProducts(page)...)
	
	log.Printf("products: %+v", len(products))
	return products
}

func getCountOfAllPages(page playwright.Page) int {
	countElement := page.Locator(".product-search-count-bottom")
	countText, err := countElement.InnerText()
	if err != nil {
		log.Printf("could not get count of all products: %v", err)
		return 0
	}

	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(countText, 2)
	if len(matches) != 2 {
		log.Printf("could not get count of all products: %v", err)
		return 0
	}

	countAll, err1 := strconv.ParseFloat(matches[1], 64)
	countPerPage, err2 := strconv.ParseFloat(matches[0], 64)
	if err1 != nil || err2 != nil || countPerPage == 0 {
		log.Printf("could not convert product counts to int: %v, %v", err1, err2)
		return 0
	}
	return int(math.Ceil(float64(countAll) / float64(countPerPage)))
}

func getProducts(page playwright.Page) []entity.Product {
	products := []entity.Product{}

	items, ok := page.Locator(".catalog-item").All()
	if ok != nil {
		log.Fatalf("could not get catalog items: %v", ok)
	}


	for _, item := range items {
		pricesBloc := item.Locator(".catalog-item__bottom")

		currentPriceElement := pricesBloc.Locator(".product-price__top").First()
		currentPriceText, err := currentPriceElement.GetAttribute("value")
		if err != nil {
			log.Printf("could not get current price: %v", err)
			continue
		}
		currentPrice := parsePrice(currentPriceText)

		oldPriceElement := pricesBloc.Locator(".product-price__bottom")
		var oldPrice float64
		count, _ := oldPriceElement.Count()
		if count > 0 {
			oldPriceText, err := oldPriceElement.First().GetAttribute("value")
			if err == nil && oldPriceText != "" {
				oldPrice = parsePrice(oldPriceText)
			} else {
				oldPrice = currentPrice
			}
		} else {
			oldPrice = currentPrice
		}

		titleElement := item.Locator(".catalog-item__title")
		title, err := titleElement.InnerText()
		if err != nil {
			log.Printf("could not get title, skipping item: %v", err)
			continue
		}

		imgElement := item.Locator(".catalog-item__img")
		imgSrc, err := imgElement.GetAttribute("src")
		if err != nil {
			log.Printf("could not get image src: %v", err)
			imgSrc = ""
		}

		product := entity.Product{
			Name:            strings.TrimSpace(title),
			DiscountedPrice: currentPrice,
			RegularPrice:    oldPrice,
			Image:           imgSrc,
		}

		log.Printf("product: %+v", product)

		products = append(products, product)

	}

	return products
}

func parsePrice(priceText string) float64 {
	re := regexp.MustCompile(`[^\d.,]`)
	cleaned := re.ReplaceAllString(priceText, "")

	cleaned = strings.Replace(cleaned, ",", ".", -1)

	price, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		log.Printf("could not parse price '%s': %v", priceText, err)
		return 0.0
	}

	return price
}
