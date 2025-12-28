package models

type Brand struct {
    ID        uint      `gorm:"primaryKey;column:brand_id"`
    Name      string    `gorm:"not null"`
    BannerURL string    `gorm:"column:banner_url"`
    Products  []Product `gorm:"foreignKey:BrandID"`
}