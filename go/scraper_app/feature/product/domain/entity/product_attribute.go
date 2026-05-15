package entity

const (
	AttributeTypeVolume = "volume"
	AttributeTypeWeight = "weight"
)

type ProductAttribute struct {
	ID    int
	Type  string
	Value string
}
