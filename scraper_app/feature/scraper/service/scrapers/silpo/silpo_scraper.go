package silpo

import (
	"log"
	"regexp"
	"sales_monitor/scraper_app/feature/scraper/utils"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"strings"
	"time"
	"github.com/playwright-community/playwright-go"
)

type SilpoScraper struct {}

func (s *SilpoScraper) GetMarketplaceName() string {
	return "Сільпо"
}

func (s *SilpoScraper) Scrape(browser playwright.Browser, url string, wordsToIgnore []string) []*entity.ScrapedProduct {
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	page.Goto(url)
	page.WaitForLoadState()

	time.Sleep(3 * time.Second)

	loadMoreButton := page.Locator(".pagination__more")

	count, _ := loadMoreButton.Count()
	if count > 0 {
		waitErr := loadMoreButton.First().WaitFor()
		if waitErr != nil {
			log.Printf("Error waiting for load more button: %v", waitErr)
		} else {
			loadMoreButton.First().ScrollIntoViewIfNeeded()
			time.Sleep(500 * time.Millisecond)

			err = loadMoreButton.First().Click()
			if err != nil {
				log.Printf("Error clicking load more button: %v", err)
			} else {
				page.WaitForLoadState()
				time.Sleep(2 * time.Second)
			}
		}
	}

	curLen := 0
	for {
		err = page.Locator(".footer").ScrollIntoViewIfNeeded()
		if err != nil {
			break
		}
		page.WaitForLoadState()

		stableCount := waitForStableElementCount(page, "silpo-products-list-item", 500*time.Millisecond, 10)
		if stableCount == -1 {
			log.Printf("Error waiting for stable element count")
			break
		}

		if curLen == stableCount {
			log.Printf("No new items, curLen: %d, len: %d", curLen, stableCount)
			break
		}
		log.Printf("New items, curLen: %d, len: %d", curLen, stableCount)

		curLen = stableCount
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

		product, err = getProductDetails(page, product)
		if err != nil {
			log.Printf("could not get product brand: %v", err)
			page.Close()
			continue
		}
		productsWithBrand = append(productsWithBrand, product)
		page.Close()
	}

	return productsWithBrand
}

func waitForStableElementCount(page playwright.Page, selector string, checkInterval time.Duration, maxChecks int) int {
	var lastCount int = -1
	stableChecks := 0

	for i := 0; i < maxChecks; i++ {
		count, err := page.Locator(selector).Count()
		if err != nil {
			log.Printf("Error counting elements: %v", err)
			return -1
		}

		if count == lastCount && lastCount >= 0 {
			stableChecks++
			if stableChecks >= 3 {
				return count
			}
		} else {
			stableChecks = 0
			lastCount = count
		}

		time.Sleep(checkInterval)
	}

	return lastCount
}

func getProducts(page playwright.Page, wordsToIgnore []string) []*entity.ScrapedProduct {
	products := []*entity.ScrapedProduct{}

	result, err := page.Evaluate(`
		() => {
			const items = document.querySelectorAll('silpo-products-list-item');
			const products = [];
			
			for (const item of items) {
				const currentPriceEl = item.querySelector('.product-card-price__displayPrice');
				if (!currentPriceEl) continue;
				
				const currentPrice = currentPriceEl.textContent?.trim() || '';
				if(!currentPrice) continue;
				
				const oldPriceEl = item.querySelector('.product-card-price__displayOldPrice');
				const oldPrice = oldPriceEl ? (oldPriceEl.textContent?.trim() || currentPrice) : currentPrice;
				
				const titleEl = item.querySelector('.product-card__title');
				if (!titleEl) continue;
				const title = titleEl.textContent?.trim() || '';
				
				const imgEl = item.querySelector('.product-card__product-img');
				let imgSrc = '';
				if (imgEl) {
					imgSrc = imgEl.getAttribute('src') || '';
				}
				
				const urlEl = item.querySelector('.product-card__link');
				const url = urlEl ? (urlEl.getAttribute('href') || '') : '';

				products.push({
					title: title,
					currentPrice: currentPrice,
					oldPrice: oldPrice,
					imgSrc: imgSrc,
					url: url,
				});
			}
			
			return products;
		}
	`)

	if err != nil {
		log.Printf("could not get products via JavaScript: %v", err)
		return products
	}

	productsData, ok := result.([]interface{})
	if !ok {
		log.Printf("unexpected result type from JavaScript")
		return products
	}

	for _, p := range productsData {
		productMap, ok := p.(map[string]interface{})
		if !ok {
			continue
		}

		title, _ := productMap["title"].(string)
		title = utils.ReplaceIgnoredWords(title, wordsToIgnore)
		currentPrice, _ := productMap["currentPrice"].(string)
		oldPrice, _ := productMap["oldPrice"].(string)
		imgSrc, _ := productMap["imgSrc"].(string)
		url, _ := productMap["url"].(string)

		if title == "" || currentPrice == "" {
			continue
		}

		product := entity.NewScrapedProduct(
			strings.TrimSpace(title),
			currentPrice,
			oldPrice,
			imgSrc,
			"https://silpo.ua" + url,
		)

		products = append(products, product)
	}

	return products
}


func getProductDetails(page playwright.Page, product *entity.ScrapedProduct) (*entity.ScrapedProduct, error) {
	title, err := page.Locator("h1").TextContent()
	if err != nil {
		log.Printf("could not get volume, weight: %v", err)
		return nil, err
	}

	re := regexp.MustCompile(`[\d|\,|\.]+[А-я-a-z|\s]+$`)
	amount := re.FindAllString(strings.TrimSpace(title), -1)
		if len(amount) == 0 || len(amount[len(amount)-1]) == 0 {
		log.Printf("could not get amount: %v", err)
		return nil, err
	}

	err = utils.ScraperSetVolumeOrWeight(amount[len(amount)-1], product)
	if err != nil {
		log.Printf("could not set volume or weight: %v", err)
		return nil, err
	}

	descriptions, err := page.Locator(".mat-expansion-panel").All()
	if err != nil {
		log.Printf("could not get brand: %v", err)
	}

	for _, item := range descriptions {	
		textContent, err := item.TextContent()
		if err != nil {
			log.Printf("could not get text content: %v", err)
			continue
		}
		if strings.Contains(textContent, "Загальна інформація") {
			description, err := item.Locator(".attributes-list_block").All()
			if err != nil {
				return nil, err
			}
			for _, attribute := range description {
				attributeTitle, _ := attribute.Locator("[data-autotestid='product-attributes-list-block-title']").TextContent()

				if  strings.TrimSpace(attributeTitle) == "Торгова марка" { 
					attributeValue, err := attribute.Locator(".attributes-list_block-value").TextContent()
					if err != nil { 
						return nil, err
					}
					product.BrandName = strings.TrimSpace(attributeValue)
					return product, nil
				}
			}
		}
		
	}
	return nil, err
}