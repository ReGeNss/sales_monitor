package fora

import (
	"log"
	"sales_monitor/scraper_app/feature/scraper/data/scraper/helper/cache"
	"sales_monitor/scraper_app/feature/scraper/data/scraper/helper/page"
	product "sales_monitor/scraper_app/shared/product/domain/entity"

	"sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/domain/gateway"
	"strings"

	"github.com/playwright-community/playwright-go"
)

type ForaScraper struct {
	Browser     playwright.Browser
	ErrorLogger gateway.ErrorLogger
}

func (s *ForaScraper) GetMarketplaceName() string {
	return "Фора"
}

func (s *ForaScraper) Scrape(url string, cachedProducts *entity.LaterScrapedProducts, wordsToIgnore []string) *entity.ScrapeResult {
	context, err := s.Browser.NewContext()
	if err != nil {
		log.Fatalf("could not create context: %v", err)
	}
	p, err := page.Open(context.Browser())
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	p.Goto(url)
	p.WaitForLoadState()

	for {
		loadMoreButton := p.Locator(".load-more-items__btn")
		count, _ := loadMoreButton.Count()
		if count > 0 {
			loadMoreButton.Click()
			p.WaitForLoadState()
		} else {
			break
		}
	}

	products := s.getProducts(p, wordsToIgnore)
	p.Close()

	productsWithBrand := []*product.ScrapedProduct{}
	newCount := 0
	for _, product := range products {
		inCache := false
		if cachedProducts != nil {
			cachedProduct, ok := (*cachedProducts)[product.URL()]
			if ok {
				inCache = true
				if cache.IsUpToDate(&cachedProduct, product) {
					continue
				}
			}
		}

		(func() {
			productURL := product.URL
			p, err = page.Open(s.Browser)
			if err != nil {
				log.Fatalf("could not create page: %v", err)
			}

			p.Goto(productURL())
			defer p.Close()
			p.WaitForLoadState()

			product, err = s.getProductBrand(p, product)
			if err != nil {
				s.logErr(p, err, gateway.ErrorContext{Context: "fora_product_brand", URL: productURL()})
				log.Printf("could not get product brand: %v", err)
				return
			}

			if err = product.Validate(); err != nil {
				return 
			}


			if !inCache {
				newCount++
			}
			productsWithBrand = append(productsWithBrand, product)
		})()
	}

	return &entity.ScrapeResult{
		Products:   productsWithBrand,
		FoundCount: len(products),
		NewCount:   newCount,
	}
}

func (s *ForaScraper) getProducts(p playwright.Page, wordsToIgnore []string) []*product.ScrapedProduct {
	products := []*product.ScrapedProduct{}

	items, err := p.Locator(".product-list-item").All()
	if err != nil {
		s.logErr(p, err, gateway.ErrorContext{Context: "fora_product_items"})
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
			s.logErr(p, err, gateway.ErrorContext{Context: "fora_product_link", Index: index, URL: productLink})
			log.Printf("could not get product link: %v", err)
		}

		currentPriceIntElement := pricesBloc.Locator(".current-integer")
		currentPriceInt, err := currentPriceIntElement.InnerText()
		if err != nil {
			s.logErr(p, err, gateway.ErrorContext{Context: "fora_current_price", Index: index, URL: productLink})
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
			s.logErr(p, err, gateway.ErrorContext{Context: "fora_title", Index: index, URL: productLink})
			log.Printf("could not get title, skipping item: %v", err)
			continue
		}

		imgElement := item.Locator(".product-list-item__image")
		imgSrc, err := imgElement.GetAttribute("src")
		if err != nil {
			s.logErr(p, err, gateway.ErrorContext{Context: "fora_image_src", Index: index, URL: productLink})
			log.Printf("could not get image src: %v", err)
			imgSrc = ""
		}

		product, err := product.CreateEmptyScrapedProduct(
			title,
			oldPrice,
			currentPrice,
			imgSrc,
			"https://fora.ua"+productLink,
			wordsToIgnore,
		)
		if err != nil {
			continue
		}

		products = append(products, product)
	}

	return products
}

func (s *ForaScraper) getProductBrand(p playwright.Page, product *product.ScrapedProduct) (*product.ScrapedProduct, error) {
	amount, err := p.Locator(".preview-product-weight").InnerText()
	if err != nil {
		s.logErr(p, err, gateway.ErrorContext{Context: "fora_amount", URL: product.URL()})
		log.Printf("could not get amount: %v", err)
		return nil, err
	}

	product.SetVolumeOrWeight(amount)

	descriptions, err := p.Locator(".product-details-column.trademark").All()
	if err != nil {
		s.logErr(p, err, gateway.ErrorContext{Context: "fora_descriptions", URL: product.URL()})
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
				s.logErr(p, err, gateway.ErrorContext{Context: "fora_description_value", URL: product.URL()})
				return nil, err
			}
			product.SetBrandName(descriptionValue)
			return product, nil
		}
	}
	return nil, err
}

func (s *ForaScraper) logErr(p playwright.Page, err error, ctx gateway.ErrorContext) {
	path := s.ErrorLogger.LogError(err, ctx)
	if path == "" || p == nil {
		return
	}
	p.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(path),
		FullPage: playwright.Bool(true),
	})
}
