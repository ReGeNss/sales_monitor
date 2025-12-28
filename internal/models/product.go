package models

type Product struct {
	ID              uint     `gorm:"primaryKey;column:product_id"`
	Name            string   `gorm:"not null; unique"`
	NameFingerprint string   `gorm:"column:name_fingerprint;not null;unique;index:idx_product_fingerprint_fulltext,type:fulltext"`
	ImageURL        string   `gorm:"column:image_url"`
	BrandID         uint     `gorm:"not null;column:brand_id"`
	Brand           Brand    `gorm:"foreignKey:BrandID"`
	CategoryID      uint     `gorm:"not null;column:category_id"`
	Category        Category `gorm:"foreignKey:CategoryID"`
	Prices          []Price  `gorm:"foreignKey:ProductID"`
}
