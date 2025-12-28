package models

type Marketplace struct {
    ID   uint   `gorm:"primaryKey;column:marketplace_id"`
    Name string `gorm:"not null"`
    URL  string `gorm:"not null"`
}