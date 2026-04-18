package models

const (
	VOLUME = "volume"
	WEIGHT = "weight"
)

type ProductAttribute struct {
	AttributeID   int       `gorm:"primaryKey;column:attribute_id;autoIncrement"`
	AttributeType string    `gorm:"notNull;column:attribute_type;type:enum('volume', 'weight')"`
	Value         string    `gorm:"notNull;column:value;type:text"`
	Products      []Product `gorm:"many2many:product_attributes;joinForeignKey:product_id;joinReferences:attribute_id"`
}

func (ProductAttribute) TableName() string {
	return "attributes"
}
