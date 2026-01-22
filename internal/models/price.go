package models

import "time"

type Price struct {
	PriceID              int                `gorm:"primaryKey;column:price_id;autoIncrement"`
	MarketplaceProductID int                `gorm:"notNull;column:marketplace_product_id"`
	MarketplaceProduct   MarketplaceProduct `gorm:"foreignKey:MarketplaceProductID;references:MarketplaceProductID"`
	RegularPrice         float64            `gorm:"notNull;column:regular_price;type:decimal(10,2)"`
	DiscountPrice        *float64           `gorm:"column:discount_price;type:decimal(10,2)"`
	CreatedAt            time.Time          `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (Price) TableName() string {
	return "prices"
}
