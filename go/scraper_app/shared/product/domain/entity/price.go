package entity

import "time"

type Price struct {
	ID                   int
	MarketplaceProductID int
	RegularPrice         float64
	SpecialPrice         *float64
	CreatedAt            time.Time
}
