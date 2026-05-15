package valueObject

import (
	"fmt"
	"net/url"
	"sales_monitor/scraper_app/feature/product/domain/exception"
)

type Url struct {
	url string
}

func (u *Url) Url() string {
	return u.url
}

func NewUrl(rawUrl string) (*Url, exception.IDomainError) {
	parsed, err := url.ParseRequestURI(rawUrl)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return nil, exception.NewDomainError(fmt.Sprintf("invalid url '%s'", rawUrl))
	}

	return &Url{url: rawUrl}, nil
}
