package models

type User struct {
	UserID   int    `gorm:"primaryKey;column:user_id;autoIncrement"`
	Login    string `gorm:"unique;notNull;column:login;type:varchar(255)"`
	Password string `gorm:"notNull;column:password;type:varchar(255)"`
	NFToken  string `gorm:"column:nf_token;type:text"`

	FavoriteBrands   []Brand   `gorm:"many2many:Favorite_Brand;joinForeignKey:user_id;joinReferences:brand_id"`
	FavoriteProducts []Product `gorm:"many2many:Favorite_Product;joinForeignKey:user_id;joinReferences:product_id"`
}
