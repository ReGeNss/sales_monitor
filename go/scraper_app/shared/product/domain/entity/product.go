package entity

import valueObject "sales_monitor/scraper_app/shared/product/domain/entity/value_object"

type Product struct {
	ID              int
	Name            string
	NameFingerprint *string
	ImageURL        valueObject.Url
	BrandID         int
	CategoryID      int
	Attributes      []*ProductAttribute
}
