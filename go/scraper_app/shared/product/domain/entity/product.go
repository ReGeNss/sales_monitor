package entity

type Product struct {
	ID              int
	Name            string
	NameFingerprint *string
	ImageURL        string
	BrandID         int
	CategoryID      int
	Attributes      []*ProductAttribute
}
