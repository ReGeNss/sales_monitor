package atb

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/domain/gateway"
	"sales_monitor/scraper_app/feature/scraper/service/dto"
	"sales_monitor/scraper_app/feature/scraper/utils"

	"strconv"
	"strings"

	"github.com/playwright-community/playwright-go"
)

type AtbScraper struct {
	Browser     playwright.Browser
	ErrorLogger gateway.ErrorLogger
}

func (s *AtbScraper) GetMarketplaceName() string {
	return "АТБ"
}

func (s *AtbScraper) Scrape(url string, cachedProducts *entity.LaterScrapedProducts) *dto.ScrapeResult {
	page, err := utils.OpenPage(s.Browser)
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	page.Goto(url)
	page.WaitForLoadState()

	countOfAllPages := s.getCountOfAllPages(page)

	var products []*dto.ScrapedProductDto

	for i := 1; i < countOfAllPages; i++ {
		products = append(products, s.getProducts(page)...)
		page.Close()

		page, err = utils.OpenPage(s.Browser)
		if err != nil {
			log.Fatalf("could not create page: %v", err)
		}
		url := fmt.Sprintf("%s?page=%d", url, i+1)
		page.Goto(url)
		page.WaitForLoadState()
	}

	products = append(products, s.getProducts(page)...)
	page.Close()

	productsWithBrand := []*dto.ScrapedProductDto{}
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
			page, err = utils.OpenPage(s.Browser)
			if err != nil {
				log.Fatalf("could not create page: %v", err)
			}
			defer page.Close()
			page.Goto(product.URL)
			page.WaitForLoadState()

			product, err = s.getProductDetails(page, product)
			if err != nil {
				s.logErr(page, err, gateway.ErrorContext{Context: "atb_product_details", URL: product.URL})
				log.Printf("could not get product brand: %v", err)
				return
			}
			if !inCache {
				newCount++
			}
			productsWithBrand = append(productsWithBrand, product)
		})()
	}

	return &dto.ScrapeResult{
		Products:   productsWithBrand,
		FoundCount: len(products),
		NewCount:   newCount,
	}
}

func (s *AtbScraper) getCountOfAllPages(page playwright.Page) int {
	countElement := page.Locator(".product-search-count-bottom")
	countText, err := countElement.InnerText()
	if err != nil {
		s.logErr(page, err, gateway.ErrorContext{Context: "atb_get_count"})
		log.Printf("could not get count of all products: %v", err)
		return 0
	}

	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(countText, 2)
	if len(matches) != 2 {
		s.logErr(page, fmt.Errorf("expected 2 matches, got %d", len(matches)), gateway.ErrorContext{Context: "atb_parse_count"})
		log.Printf("could not get count of all products: matches=%d", len(matches))
		return 0
	}

	countAll, err1 := strconv.ParseFloat(matches[1], 64)
	countPerPage, err2 := strconv.ParseFloat(matches[0], 64)
	if err1 != nil || err2 != nil || countPerPage == 0 {
		s.logErr(page, err1, gateway.ErrorContext{Context: "atb_convert_count"})
		log.Printf("could not convert product counts to int: %v, %v", err1, err2)
		return 0
	}
	return int(math.Ceil(float64(countAll) / float64(countPerPage)))
}

func (s *AtbScraper) getProducts(page playwright.Page) []*dto.ScrapedProductDto {
	products := []*dto.ScrapedProductDto{}

	items, ok := page.Locator(".catalog-item").All()
	if ok != nil {
		s.logErr(page, ok, gateway.ErrorContext{Context: "atb_catalog_items"})
		log.Fatalf("could not get catalog items: %v", ok)
	}

	for index, item := range items {
		pricesBloc := item.Locator(".catalog-item__bottom")

		productLink, err := item.Locator(".catalog-item__photo-link").GetAttribute("href")
		if err != nil {
			s.logErr(page, err, gateway.ErrorContext{Context: "atb_product_link_index", Index: index})
			log.Printf("could not get product link: %v", err)
		}

		currentPriceElement := pricesBloc.Locator(".product-price__top").First()
		currentPriceText, err := currentPriceElement.GetAttribute("value")
		if err != nil {
			s.logErr(page, err, gateway.ErrorContext{Context: "atb_current_price", URL: productLink, Index: index})
			log.Printf("could not get current price: %v", err)
			continue
		}
		currentPrice := currentPriceText

		oldPriceElement := pricesBloc.Locator(".product-price__bottom")
		var oldPrice string
		count, _ := oldPriceElement.Count()
		if count > 0 {
			oldPrice, _ = oldPriceElement.First().GetAttribute("value")
		} else {
			oldPrice = currentPrice
		}

		titleElement := item.Locator(".catalog-item__title")
		title, err := titleElement.InnerText()
		if err != nil {
			s.logErr(page, err, gateway.ErrorContext{Context: "atb_title", Index: index, URL: productLink})
			log.Printf("could not get title, skipping item: %v", err)
			continue
		}

		imgElement := item.Locator(".catalog-item__img")
		imgSrc, err := imgElement.GetAttribute("src")
		if err != nil {
			s.logErr(page, err, gateway.ErrorContext{Context: "atb_image_src", Index: index, URL: productLink})
			log.Printf("could not get image src: %v", err)
			imgSrc = ""
		}

		product := dto.CreateScrapedProductDto(
			strings.TrimSpace(title),
			oldPrice,
			currentPrice,
			imgSrc,
			"https://www.atbmarket.com"+productLink,
		)

		products = append(products, product)

	}

	return products
}

func (s *AtbScraper) getProductDetails(page playwright.Page, product *dto.ScrapedProductDto) (*dto.ScrapedProductDto, error) {
	brandElement, err := page.Locator(".product-characteristics__item").All()
	if err != nil {
		s.logErr(page, err, gateway.ErrorContext{Context: "atb_brand_elements", URL: product.URL})
		log.Printf("could not get brand: %v", err)
		return nil, err
	}
	for _, item := range brandElement {
		elementTitle, err := item.Locator(".product-characteristics__name").InnerText()
		if err != nil {
			continue
		}

		if elementTitle == "Торгова марка" {
			brandName, err := s.getProductAttributeValue(page, item, product.URL)
			if err != nil {
				s.logErr(page, err, gateway.ErrorContext{Context: "atb_brand_name", URL: product.URL})
				log.Printf("could not get brand name: %v", err)
				return nil, err
			}
			product.BrandName = brandName
		}

		if elementTitle == "Об’єм" {
			volume, err := s.getProductAttributeValue(page, item, product.URL)
			if err == nil {
				product.ScraperSetVolumeOrWeight(volume)
			}
		}

		if elementTitle == "Вага" {
			weight, err := s.getProductAttributeValue(page, item, product.URL)
			if err == nil {
				product.ScraperSetVolumeOrWeight(weight)
			}
		}
	}
	return product, nil
}

func (s *AtbScraper) getProductAttributeValue(page playwright.Page, item playwright.Locator, url string) (string, error) {
	volumeElement, err := item.Locator(".product-characteristics__value").InnerText()
	if err != nil {
		s.logErr(page, err, gateway.ErrorContext{Context: "atb_attribute_value", URL: url})
		log.Printf("could not get volume: %v", err)
		return "", err
	}
	return volumeElement, nil
}

func (s *AtbScraper) logErr(page playwright.Page, err error, ctx gateway.ErrorContext) {
	path := s.ErrorLogger.LogError(err, ctx)
	if path == "" || page == nil {
		return
	}
	page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(path),
		FullPage: playwright.Bool(true),
	})
}
