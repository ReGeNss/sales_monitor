package scraper

import (
	"log"
	"sales_monitor/scraper_app/feature/scraper/data/scraper/atb"
	"sales_monitor/scraper_app/feature/scraper/data/scraper/fora"
	"sales_monitor/scraper_app/feature/scraper/data/scraper/silpo"
	"sales_monitor/scraper_app/feature/scraper/domain/exception"
	"sales_monitor/scraper_app/feature/scraper/domain/gateway"

	"github.com/playwright-community/playwright-go"
)

const (
	ATB   = "atb"
	FORA  = "fora"
	SILPO = "silpo"
)

type playwrightScraperFactory struct {
	pw          *playwright.Playwright
	browser     playwright.Browser
	errorLogger gateway.ErrorLogger
}

func NewScraperFactory(errorLogger gateway.ErrorLogger) (gateway.ScraperFactory, exception.IDomainError) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, exception.NewDomainError("couldn't start scraper: " + err.Error())
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args: []string{
			"--disable-blink-features=AutomationControlled",
			"--disable-dev-shm-usage",
		},
	})
	if err != nil {
		pw.Stop()
		return nil, exception.NewDomainError("couldn't start browser: " + err.Error())
	}

	return &playwrightScraperFactory{
		pw:          pw,
		browser:     browser,
		errorLogger: errorLogger,
	}, nil
}

func (f *playwrightScraperFactory) Get(shopID string) (gateway.Scraper, exception.IDomainError) {
	switch shopID {
	case ATB:
		return &atb.AtbScraper{Browser: f.browser, ErrorLogger: f.errorLogger}, nil
	case FORA:
		return &fora.ForaScraper{Browser: f.browser, ErrorLogger: f.errorLogger}, nil
	case SILPO:
		return &silpo.SilpoScraper{Browser: f.browser, ErrorLogger: f.errorLogger}, nil
	default:
		return nil, exception.NewDomainError("Unknown shop id: " + shopID)
	}
}

func (f *playwrightScraperFactory) Close() {
	if err := f.browser.Close(); err != nil {
		log.Printf("could not close browser: %v", err)
	}
	if err := f.pw.Stop(); err != nil {
		log.Printf("could not stop playwright: %v", err)
	}
}
