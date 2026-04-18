package scrapers

import (
	"fmt"

	scraper_entity "sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/atb"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/fora"
	"sales_monitor/scraper_app/feature/scraper/service/scrapers/silpo"
)

func GetScraperByShopName(shopID string) (scraper_entity.Scraper, error) {
	switch shopID {
	case "atb":
		return &atb.AtbScraper{}, nil
	case "fora":
		return &fora.ForaScraper{}, nil
	case "silpo":
		return &silpo.SilpoScraper{}, nil
	default:
		return nil, fmt.Errorf("unknown shop_id %q", shopID)
	}
}
