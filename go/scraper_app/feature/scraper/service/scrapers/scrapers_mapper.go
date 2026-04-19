package scrapers

import (
	"fmt"

	"sales_monitor/scraper_app/feature/scraper/domain/gateway"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/atb"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/fora"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/silpo"

	"github.com/playwright-community/playwright-go"
)

const (
	ATB   = "atb"
	FORA  = "fora"
	SILPO = "silpo"
)

func GetScraperByShopName(shopID string, browser playwright.Browser, errorLogger gateway.ErrorLogger) (Scraper, error) {
	switch shopID {
	case ATB:
		return &atb.AtbScraper{Browser: browser, ErrorLogger: errorLogger}, nil
	case FORA:
		return &fora.ForaScraper{Browser: browser, ErrorLogger: errorLogger}, nil
	case SILPO:
		return &silpo.SilpoScraper{Browser: browser, ErrorLogger: errorLogger}, nil
	default:
		return nil, fmt.Errorf("unknown shop_id %q", shopID)
	}
}
