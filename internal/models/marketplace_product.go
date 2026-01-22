package models

type MarketplaceProduct struct {
	MarketplaceProductID int         `gorm:"primaryKey;column:marketplace_product_id;autoIncrement"`
	MarketplaceID        int         `gorm:"notNull;column:marketplace_id"`
	ProductID            int         `gorm:"notNull;column:product_id"`
	URL                  string      `gorm:"notNull;column:url;type:text"`
	Marketplace          Marketplace `gorm:"foreignKey:MarketplaceID;references:MarketplaceID"`
	Product              Product     `gorm:"foreignKey:ProductID;references:ProductID"`
	Prices               []Price     `gorm:"foreignKey:MarketplaceProductID"`
}

func (MarketplaceProduct) TableName() string {
	return "marketplace_products"
}
