package models

const (
	VOLUME = "volume"
	WEIGHT = "weight"
)

type ProductAttribute struct {
	ProductAttributeID int `gorm:"primaryKey;column:product_attribute_id;autoIncrement"`
	ProductID          int `gorm:"notNull;column:product_id"`
	AttributeType      string `gorm:"notNull;column:attribute_type;type:enum('volume', 'weight')"`
	Value              string `gorm:"notNull;column:value;type:text"`
}
