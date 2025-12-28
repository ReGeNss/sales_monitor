package models

import "time"

type Price struct {
	ID              uint        `gorm:"primaryKey;column:price_id"`
	ProductID       uint        `gorm:"not null;column:product_id"`
	Product         Product     `gorm:"foreignKey:ProductID"`
	MarketplaceID   uint        `gorm:"not null;column:marketplace_id"`
	Marketplace     Marketplace `gorm:"foreignKey:MarketplaceID"`
	RegularPrice    float64     `gorm:"column:regular_price;type:decimal(10,2);not null"`
	DiscountPrice   float64     `gorm:"column:discount_price;type:decimal(10,2)"`
	URL             string      `gorm:"not null"`
	IsOnSale        bool        `gorm:"not null;column:is_on_sale;default:false"`
	DiscountPercent int         `gorm:"not null;column:discount_percent;default:0"`
	CreatedAt       time.Time   `gorm:"autoCreateTime"`
}
