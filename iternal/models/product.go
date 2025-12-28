package models

type Product struct {
    ID         uint      `gorm:"primaryKey;column:product_id"`
    Name       string    `gorm:"not null"`
    ImageURL   string    `gorm:"column:image_url"`
    BrandID    uint      `gorm:"column:brand_id"`
    Brand      Brand     `gorm:"foreignKey:BrandID"`
    CategoryID uint      `gorm:"column:category_id"`
    Category   Category  `gorm:"foreignKey:CategoryID"`
    Prices     []Price   `gorm:"foreignKey:ProductID"`
}