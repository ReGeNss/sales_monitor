package models

type User struct {
	ID                uint      `gorm:"primaryKey;column:user_id"`
	Login             string    `gorm:"unique;not null"`
	Password          string    `gorm:"not null"`
	NFToken           string    `gorm:"column:nf_token"`
	FavoritesProducts []Product `gorm:"many2many:favorite_product;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:ProductID"`
	FavoritesBrands   []Brand   `gorm:"many2many:favorite_brand;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:BrandID"`
}
