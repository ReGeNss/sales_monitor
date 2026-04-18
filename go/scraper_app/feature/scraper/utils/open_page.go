package utils

import "github.com/playwright-community/playwright-go"

func OpenPage(browser playwright.Browser) (playwright.Page, error) {
	return browser.NewPage(
		playwright.BrowserNewPageOptions{
			UserAgent: playwright.String("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		},
	)
}
