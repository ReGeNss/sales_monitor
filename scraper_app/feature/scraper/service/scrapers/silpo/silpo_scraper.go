package silpo

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"regexp"
	scraper_config "sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/utils"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"strings"
	"time"
)

type SilpoScraper struct{}

func (s *SilpoScraper) GetMarketplaceName() string {
	return "Сільпо"
}

func (s *SilpoScraper) Scrape(browser playwright.Browser, url string, wordsToIgnore []string, cachedProducts *scraper_config.LaterScrapedProducts) *scraper_config.ScrapeResult {
	page, err := utils.OpenPage(browser)
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	page.Goto(url)
	page.WaitForLoadState()

	page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("screenshot.png"),
		FullPage: playwright.Bool(true),
	})

	time.Sleep(3 * time.Second)

	loadMoreButton := page.Locator(".pagination__more")

	count, _ := loadMoreButton.Count()
	if count > 0 {
		waitErr := loadMoreButton.First().WaitFor()
		if waitErr != nil {
			utils.SaveScreenshotOnError(page, waitErr, "silpo_load_more_wait")
			log.Printf("Error waiting for load more button: %v", waitErr)
		} else {
			loadMoreButton.First().ScrollIntoViewIfNeeded()
			time.Sleep(500 * time.Millisecond)

			err = loadMoreButton.First().Click()
			if err != nil {
				utils.SaveScreenshotOnError(page, err, "silpo_load_more_click")
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
			utils.SaveScreenshotOnError(page, fmt.Errorf("stable element count failed"), "silpo_stable_count")
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
			defer page.Close()
			page.Goto(product.URL)
			page.WaitForLoadState()

			product, err = getProductDetails(page, product)
			if err != nil {
				utils.SaveScreenshotOnError(page, err, "silpo_product_details")
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

func waitForStableElementCount(page playwright.Page, selector string, checkInterval time.Duration, maxChecks int) int {
	var lastCount int = -1
	stableChecks := 0

	for i := 0; i < maxChecks; i++ {
		count, err := page.Locator(selector).Count()
		if err != nil {
			utils.SaveScreenshotOnError(page, err, "silpo_count_elements")
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
		utils.SaveScreenshotOnError(page, err, "silpo_js_products")
		log.Printf("could not get products via JavaScript: %v", err)
		return products
	}

	productsData, ok := result.([]interface{})
	if !ok {
		utils.SaveScreenshotOnError(page, fmt.Errorf("unexpected result type"), "silpo_js_result")
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
			"https://silpo.ua"+url,
		)

		products = append(products, product)
	}

	return products
}

func getProductDetails(page playwright.Page, product *entity.ScrapedProduct) (*entity.ScrapedProduct, error) {
	title, err := page.Locator("h1").TextContent()
	if err != nil {
		utils.SaveScreenshotOnError(page, err, "silpo_title")
		log.Printf("could not get volume, weight: %v", err)
		return nil, err
	}

	re := regexp.MustCompile(`[\d|\,|\.]+[А-я-a-z|\s]+$`)
	amount := re.FindAllString(strings.TrimSpace(title), -1)
	if len(amount) == 0 || len(amount[len(amount)-1]) == 0 {
		utils.SaveScreenshotOnError(page, err, "silpo_amount")
		log.Printf("could not get amount: %v", err)
		return nil, err
	}

	err = utils.ScraperSetVolumeOrWeight(amount[len(amount)-1], product)
	if err != nil {
		utils.SaveScreenshotOnError(page, err, "silpo_volume_weight")
		log.Printf("could not set volume or weight: %v", err)
		return nil, err
	}

	descriptions, err := page.Locator(".mat-expansion-panel").All()
	if err != nil {
		utils.SaveScreenshotOnError(page, err, "silpo_descriptions")
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
				utils.SaveScreenshotOnError(page, err, "silpo_attributes")
				return nil, err
			}
			for _, attribute := range description {
				attributeTitle, _ := attribute.Locator("[data-autotestid='product-attributes-list-block-title']").TextContent()

				if strings.TrimSpace(attributeTitle) == "Торгова марка" {
					attributeValue, err := attribute.Locator(".attributes-list_block-value").TextContent()
					if err != nil {
						utils.SaveScreenshotOnError(page, err, "silpo_brand_value")
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
