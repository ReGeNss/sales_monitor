package fora

import (
	"log"
	"sales_monitor/scraper_app/feature/scraper/utils"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"strings"
	"github.com/playwright-community/playwright-go"
)

type ForaScraper struct {}

func (s *ForaScraper) GetMarketplaceName() string {
	return "Фора"
}

func (s *ForaScraper) Scrape(browser playwright.Browser, url string, wordsToIgnore []string) []*entity.ScrapedProduct {
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

	products := getProducts(page, wordsToIgnore)

	productsWithBrand := []*entity.ScrapedProduct{}
	for _, product := range products {
		page, err = browser.NewPage()
		if err != nil {
			log.Fatalf("could not create page: %v", err)
		}

		page.Goto(product.URL)
		page.WaitForLoadState()

		product, err = getProductBrand(page, product)
		if err != nil {
			log.Printf("could not get product brand: %v", err)
			continue
		}
		productsWithBrand = append(productsWithBrand, product)
		page.Close()
	}

	return productsWithBrand
}

func getProducts(page playwright.Page, wordsToIgnore []string) []*entity.ScrapedProduct {
	products := []*entity.ScrapedProduct{}

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

		title = utils.ReplaceIgnoredWords(title, wordsToIgnore)

		imgElement := item.Locator(".product-list-item__image")
		imgSrc, err := imgElement.GetAttribute("src")
		if err != nil {
			log.Printf("could not get image src: %v", err)
			imgSrc = ""
		}

		productLink, err := item.Locator(".image-content-wrapper").GetAttribute("href")
		if err != nil {
			log.Printf("could not get product link: %v", err)
		}

		product := entity.NewScrapedProduct(
			strings.TrimSpace(title),
			currentPrice,
			oldPrice,
			imgSrc,
			"https://fora.ua"+productLink,
		)

		products = append(products, product)
	}

	return products
}

func getProductBrand(page playwright.Page, product *entity.ScrapedProduct) (*entity.ScrapedProduct, error) {
	amount, err := page.Locator(".preview-product-weight").InnerText()

	if err != nil {
		log.Printf("could not get amount: %v", err)
		return nil, err
	}

	if strings.Contains(amount, "л") {
		utils.ScraperSetVolumeOrWeight(amount, product)
	} else {
		utils.ScraperSetVolumeOrWeight(amount, product)
	}

	descriptions, err := page.Locator(".product-details-column.trademark").All()
	
	if err != nil {
		log.Printf("could not get descriptions: %v", err)
		return nil, err
	}

	for _, description := range descriptions {
		descriptionLabel, err := description.Locator(".product-details-label").TextContent()
		if err != nil {
			log.Printf("could not get description label: %v", err)
			continue
		}
		if strings.TrimSpace(descriptionLabel) == "Торгова марка" {
			descriptionValue, err := description.Locator(".product-details-value").TextContent()
			if err != nil {
				return nil, err
			}
			product.BrandName = strings.TrimSpace(descriptionValue)
			return product, nil
		}
	}
	return nil, err
}