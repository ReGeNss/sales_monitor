package models

type NotificationTask struct {
	BrandID  int
	BrandName string
	Products []NotificationProduct
}

type NotificationProduct struct {
	ID   int
	Name string
	ImageURL string
}