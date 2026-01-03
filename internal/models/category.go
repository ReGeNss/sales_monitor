package models

type Category struct {
	CategoryID int    `gorm:"primaryKey;column:category_id;autoIncrement"`
	Name       string `gorm:"unique;notNull;column:name;type:varchar(255)"`
}
