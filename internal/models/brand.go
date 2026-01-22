package models

type Brand struct {
	BrandID   int    `gorm:"primaryKey;column:brand_id;autoIncrement"`
	Name      string `gorm:"unique;notNull;column:name;type:varchar(255)"`
	BannerURL string `gorm:"column:banner_url;type:text"`
}


func (Brand) TableName() string {
	return "brands"
}