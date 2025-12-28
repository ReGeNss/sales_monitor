package models

type Category struct {
    ID        uint      `gorm:"primaryKey;column:category_id"`
    Name      string    `gorm:"not null"`
    FieldName string    `gorm:"column:fieldname"`
    Products  []Product `gorm:"foreignKey:CategoryID"`
}