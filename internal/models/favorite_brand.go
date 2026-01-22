package models

type FavoriteBrand struct {
	FavoriteBrandID int   `gorm:"primaryKey;column:favorite_brand_id;autoIncrement"`
	UserID          int   `gorm:"notNull;column:user_id"`
	BrandID         int   `gorm:"notNull;column:brand_id"`
	User            User  `gorm:"foreignKey:UserID;references:user_id"`
	Brand           Brand `gorm:"foreignKey:BrandID;references:brand_id"`
}

func (FavoriteBrand) TableName() string {
	return "favorite_brands"
}
