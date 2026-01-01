package atb

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"sales_monitor/scraper_app/feature/scraper/utils"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"strconv"
	"strings"

	"github.com/playwright-community/playwright-go"
)

func AtbScraper(browser playwright.Browser, url string) []*entity.ScrapedProduct {
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	page.Goto(url)
	page.WaitForLoadState()

	countOfAllPages := getCountOfAllPages(page)

	var products []*entity.ScrapedProduct

	for i := 1; i < countOfAllPages; i++ {
		products = append(products, getProducts(page)...)
		page.Close()

		page, err = browser.NewPage()
		if err != nil {
			log.Fatalf("could not create page: %v", err)
		}
		url := fmt.Sprintf("%s?page=%d", url, i+1)
		page.Goto(url)
		page.WaitForLoadState()
	}
	
	products = append(products, getProducts(page)...)
	page.Close()

	productsWithBrand := []*entity.ScrapedProduct{}
	for _, product := range products {
		page, err = browser.NewPage()
		if err != nil {
			log.Fatalf("could not create page: %v", err)
		}

		page.Goto(product.URL)
		page.WaitForLoadState()

		product, err = getProductDetails(page, product)
		if err != nil {
			log.Printf("could not get product brand: %v", err)
			continue
		}
		productsWithBrand = append(productsWithBrand, product)
		page.Close()
	}

	return productsWithBrand
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

func getProducts(page playwright.Page) []*entity.ScrapedProduct {
	products := []*entity.ScrapedProduct{}

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
			log.Printf("could not get title, skipping item: %v", err)
			continue
		}

		imgElement := item.Locator(".catalog-item__img")
		imgSrc, err := imgElement.GetAttribute("src")
		if err != nil {
			log.Printf("could not get image src: %v", err)
			imgSrc = ""
		}

		productLink, err := item.Locator(".catalog-item__photo-link").GetAttribute("href")
		if err != nil {
			log.Printf("could not get product link: %v", err)
		}

		product := entity.NewScrapedProduct(
			strings.TrimSpace(title),
			currentPrice,
			oldPrice,
			imgSrc,
			"https://www.atbmarket.com" + productLink,
		)

		products = append(products, product)

	}

	return products
}


func getProductDetails(page playwright.Page, product *entity.ScrapedProduct) (*entity.ScrapedProduct, error) {
	brandElement, err := page.Locator(".product-characteristics__item").All()
	if err != nil {
		log.Printf("could not get brand: %v", err)
		return nil, err
	}
	for _, item := range brandElement {
		elementTitle, err := item.Locator(".product-characteristics__name").InnerText()
		if err != nil {
			continue
		}

		if elementTitle == "Торгова марка" {
			brandName, err := getProductAttributeValue(item)
			if err != nil {
				log.Printf("could not get brand name: %v", err)
				return nil, err
			}
			product.BrandName = brandName
		}

		if elementTitle == "Об'єм" {
			volume, err := getProductAttributeValue(item)
			if err == nil {
				utils.ScraperSetVolumeOrWeight(volume, product)
			}
		}

		if elementTitle == "Вага" {
			weight, err := getProductAttributeValue(item)
			if err == nil {
				utils.ScraperSetVolumeOrWeight(weight, product)
			}
		}
	}
	return product, nil
}

func getProductAttributeValue(item playwright.Locator) (string, error) {
	volumeElement, err := item.Locator(".product-characteristics__value").InnerText()
	if err != nil {
		log.Printf("could not get volume: %v", err)
		return "", err
	}
	return volumeElement, nil
}