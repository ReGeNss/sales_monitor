package models
    
type Marketplace struct {
    ID   uint   `gorm:"primaryKey;column:marketplace_id"`
    Name string `gorm:"not null;unique"`
    URL  string `gorm:"not null;unique"`
}