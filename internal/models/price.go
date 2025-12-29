package models

import "time"

type Price struct {
	PriceID         int       `gorm:"primaryKey;column:price_id;autoIncrement"`
	ProductID       int       `gorm:"notNull;column:product_id"`
	RegularPrice    float64   `gorm:"notNull;column:regular_price;type:decimal(10,2)"`
	DiscountPrice   *float64  `gorm:"column:discount_price;type:decimal(10,2)"`
	MarketplaceID   int       `gorm:"notNull;column:marketplace_id"`
	URL             string    `gorm:"notNull;column:url;type:text"`
	IsOnSale        bool      `gorm:"column:is_on_sale;default:false"`
	DiscountPercent int       `gorm:"column:discount_percent;default:0"`
	CreatedAt       time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (Price) TableName() string {
	return "Price"
}
