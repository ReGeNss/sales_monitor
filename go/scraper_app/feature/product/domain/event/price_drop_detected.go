package event

type PriceDropDetected struct {
	BrandID   int
	BrandName string
	Products  []DroppedProduct
}

type DroppedProduct struct {
	ID       int
	Name     string
	ImageURL string
}
