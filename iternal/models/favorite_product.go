package models

type FavoriteProduct struct {
    ID uint `gorm:"primaryKey;column:favorite_product_id"`
    UserID     uint `gorm:"column:user_id"`
    ProductID  uint `gorm:"column:product_id"`
}