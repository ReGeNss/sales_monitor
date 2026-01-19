package models

type Product struct {
	ProductID       int                `gorm:"primaryKey;column:product_id;autoIncrement"`
	NameFingerprint *string            `gorm:"unique;column:name_fingerprint;type:varchar(255)"`
	BrandID         int                `gorm:"notNull;column:brand_id"`
	Name            string             `gorm:"unique;notNull;column:name;type:varchar(255)"`
	CategoryID      int                `gorm:"notNull;column:category_id"`
	ImageURL        string             `gorm:"column:image_url;type:text"`
	Attributes      []ProductAttribute `gorm:"many2many:product_attributes;joinForeignKey:product_id;joinReferences:attribute_id"`
	Prices          []Price            `gorm:"foreignKey:ProductID"`
}

func (Product) TableName() string {
	return "Product"
}
