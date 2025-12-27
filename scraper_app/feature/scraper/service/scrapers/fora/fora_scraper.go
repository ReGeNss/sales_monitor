package fora

import (
	"log"
	"sales_monitor/scraper/shared/product/domain/entity"
	"strings"
	"github.com/playwright-community/playwright-go"
)

func ForaScraper(browser playwright.Browser, url string) []*entity.Product {
	context, err := browser.NewContext()
	if err != nil {
		log.Fatalf("could not create context: %v", err)
	}
	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	page.Goto(url)
	page.WaitForLoadState()

	for {
		loadMoreButton := page.Locator(".load-more-items__btn")
		count, _ := loadMoreButton.Count()
		if count > 0 {
			loadMoreButton.Click()
			page.WaitForLoadState()
		} else {
			break
		}
	}

	products := getProducts(page)

	return products

}

func getProducts(page playwright.Page) []*entity.Product {
	products := []*entity.Product{}

	items, err := page.Locator(".product-list-item").All()
	if err != nil {
		log.Printf("could not get product items: %v", err)
		return products
	}

	for _, item := range items {
		pricesBloc := item.Locator(".product-price-container")
		count, _ := pricesBloc.Count()
		if count == 0 {
			continue
		}

		currentPriceIntElement := pricesBloc.Locator(".current-integer")
		currentPriceInt, err := currentPriceIntElement.InnerText()
		if err != nil {
			log.Printf("could not get current price integer: %v", err)
			continue
		}

		currentPriceFractionElement := pricesBloc.Locator(".current-fraction")
		var currentPriceFraction string
		countFraction, _ := currentPriceFractionElement.Count()
		if countFraction > 0 {
			currentPriceFraction, _ = currentPriceFractionElement.First().InnerText()
		}

		currentPrice := strings.TrimSpace(currentPriceInt) + "." + strings.TrimSpace(currentPriceFraction)

		oldPriceElement := pricesBloc.Locator(".old-integer")
		var oldPrice string
		countOld, _ := oldPriceElement.Count()
		if countOld > 0 {
			oldPrice, _ = oldPriceElement.First().InnerText()
		} else {
			oldPrice = currentPrice
		}

		titleElement := item.Locator(".product-title")
		countTitle, _ := titleElement.Count()
		if countTitle == 0 {
			continue
		}
		title, err := titleElement.InnerText()
		if err != nil {
			log.Printf("could not get title, skipping item: %v", err)
			continue
		}

		imgElement := item.Locator(".product-list-item__image")
		imgSrc, err := imgElement.GetAttribute("src")
		if err != nil {
			log.Printf("could not get image src: %v", err)
			imgSrc = ""
		}

		product := entity.NewProduct(
			strings.TrimSpace(title),
			currentPrice,
			oldPrice,
			imgSrc,
		)

		products = append(products, product)
	}

	return products
}
