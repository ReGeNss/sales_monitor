package fora

import (
	"log"
	"sales_monitor/scraper_app/feature/scraper/utils"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	scraper_config "sales_monitor/scraper_app/feature/scraper/domain/entity"
	"strings"
	"github.com/playwright-community/playwright-go"
)

type ForaScraper struct {}

func (s *ForaScraper) GetMarketplaceName() string {
	return "Фора"
}

func (s *ForaScraper) Scrape(browser playwright.Browser, url string, wordsToIgnore []string, cachedProducts *scraper_config.LaterScrapedProducts) *scraper_config.ScrapeResult {
	context, err := browser.NewContext()
	if err != nil {
		log.Fatalf("could not create context: %v", err)
	}
	page, err := utils.OpenPage(context.Browser())
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
	page.Close()
	
	productsWithBrand := []*entity.ScrapedProduct{}
	newCount := 0
	for _, product := range products {
		inCache := false
		if cachedProducts != nil {
			cachedProduct, ok := (*cachedProducts)[product.URL]
			if ok {
				inCache = true
				if utils.CheckForProductUpdate(&cachedProduct, product) {
					continue
				}
			}
		}

		(func() {
			page, err = utils.OpenPage(browser)
			if err != nil {
				log.Fatalf("could not create page: %v", err)
			}

			page.Goto(product.URL)
			defer page.Close()
			page.WaitForLoadState()

			product, err = getProductBrand(page, product)
			if err != nil {
				utils.SaveScreenshotOnError(page, err, utils.ErrorContext{Context: "fora_product_brand", URL: product.URL})
				log.Printf("could not get product brand: %v", err)
				return
			}
			if !inCache {
				newCount++
			}
			productsWithBrand = append(productsWithBrand, product)
		})()
	}

	return &scraper_config.ScrapeResult{
		Products:   productsWithBrand,
		FoundCount: len(products),
		NewCount:   newCount,
	}
}

func getProducts(page playwright.Page, wordsToIgnore []string) []*entity.ScrapedProduct {
	products := []*entity.ScrapedProduct{}

	items, err := page.Locator(".product-list-item").All()
	if err != nil {
		utils.SaveScreenshotOnError(page, err, utils.ErrorContext{Context: "fora_product_items"})
		log.Printf("could not get product items: %v", err)
		return products
	}

	for index, item := range items {
		pricesBloc := item.Locator(".product-price-container")
		count, _ := pricesBloc.Count()
		if count == 0 {
			continue
		}

		productLink, err := item.Locator(".image-content-wrapper").GetAttribute("href")
		if err != nil {
			utils.SaveScreenshotOnError(page, err, utils.ErrorContext{Context: "fora_product_link", Index: index, URL: productLink})
			log.Printf("could not get product link: %v", err)
		}

		currentPriceIntElement := pricesBloc.Locator(".current-integer")
		currentPriceInt, err := currentPriceIntElement.InnerText()
		if err != nil {
			utils.SaveScreenshotOnError(page, err, utils.ErrorContext{Context: "fora_current_price", Index: index, URL: productLink})
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
			utils.SaveScreenshotOnError(page, err, utils.ErrorContext{Context: "fora_title", Index: index, URL: productLink})
			log.Printf("could not get title, skipping item: %v", err)
			continue
		}

		title = utils.ReplaceIgnoredWords(title, wordsToIgnore)

		imgElement := item.Locator(".product-list-item__image")
		imgSrc, err := imgElement.GetAttribute("src")
		if err != nil {
			utils.SaveScreenshotOnError(page, err, utils.ErrorContext{Context: "fora_image_src", Index: index, URL: productLink})
			log.Printf("could not get image src: %v", err)
			imgSrc = ""
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
		utils.SaveScreenshotOnError(page, err, utils.ErrorContext{Context: "fora_amount", URL: product.URL})
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
		utils.SaveScreenshotOnError(page, err, utils.ErrorContext{Context: "fora_descriptions", URL: product.URL})
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
				utils.SaveScreenshotOnError(page, err, utils.ErrorContext{Context: "fora_description_value", URL: product.URL})
				return nil, err
			}
			product.BrandName = strings.TrimSpace(descriptionValue)
			return product, nil
		}
	}
	return nil, err
}