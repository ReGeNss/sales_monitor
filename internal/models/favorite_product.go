package models

type FavoriteProduct struct {
	FavoriteProductID int     `gorm:"primaryKey;column:favorite_product_id;autoIncrement"`
	UserID            int     `gorm:"notNull;column:user_id"`
	ProductID         int     `gorm:"notNull;column:product_id"`
	User              User    `gorm:"foreignKey:UserID;references:user_id"`
	Product           Product `gorm:"foreignKey:ProductID;references:product_id"`
}

func (FavoriteProduct) TableName() string {
	return "Favorite_Product"
}
