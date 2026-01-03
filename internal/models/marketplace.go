package models

type Marketplace struct {
	MarketplaceID int    `gorm:"primaryKey;column:marketplace_id;autoIncrement"`
	Name          string `gorm:"unique;notNull;column:name;type:varchar(255)"`
	URL           string `gorm:"unique;notNull;column:url;type:text"`
}
