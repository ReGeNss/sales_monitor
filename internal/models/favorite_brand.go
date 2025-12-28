package models

type FavoriteBrand struct {
    ID uint `gorm:"primaryKey;column:favorite_brand_id"`
    UserID uint `gorm:"column:user_id;not null"`
    BrandID uint `gorm:"column:brand_id;not null"`
}