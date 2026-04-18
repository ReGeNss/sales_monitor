package main

import (
	"fmt"
	"log"
	"sales_monitor/internal/models"

	"firebase.google.com/go/v4/messaging"
	"gorm.io/gorm"
)

func sendNotification(db *gorm.DB, messagingClient *messaging.Client, notificationTask models.NotificationTask) {
	tokens, err := getTokensWithFavoriteBrand(db, notificationTask.BrandID)
	if err != nil {
		log.Printf("error getting users with favorite brand: %v", err)
	}

	batchResponse, err := messagingClient.SendEachForMulticast(ctx, &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title:    fmt.Sprintf("Зʼявились нові знижки для бренду %s", notificationTask.BrandName),
			Body:     "Перейдіть в застосунок, щоб переглянути їх",
			ImageURL: notificationTask.Products[0].ImageURL,
		},
	})

	handleBatchResponse(batchResponse)

	if err != nil {
		log.Printf("error sending notification to favorite brand: %v", err)
	}

	for _, product := range notificationTask.Products {
		tokens, err = getTokensWithFavoriteProduct(db, product.ProductID, notificationTask.BrandID)
		if err != nil {
			log.Printf("error getting users with favorite product: %v", err)
		}

		batchResponse, _ := messagingClient.SendEachForMulticast(ctx, &messaging.MulticastMessage{
			Tokens: tokens,
			Notification: &messaging.Notification{
				Title:    fmt.Sprintf("Зʼявилась нова знижка для %s", product.Name),
				Body:     "Перейдіть в застосунок, щоб переглянути її",
				ImageURL: product.ImageURL,
			},
		})

		handleBatchResponse(batchResponse)
	}
}

func getTokensWithFavoriteBrand(db *gorm.DB, brandID int) ([]string, error) {
	var tokens []string
	err := db.Model(&models.User{}).
		Joins("favorite_brand fb on fb.user_id = users.user_id").
		Where("fb.brand_id = ?", brandID).
		Pluck("users.nf_token", &tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func getTokensWithFavoriteProduct(db *gorm.DB, productID int, brandId int) ([]string, error) {
	var tokens []string
	err := db.Model(&models.User{}).
		Joins("favorite_product fp on fp.user_id = users.user_id").
		Joins("LEFT JOIN favorite_brand fb on fb.user_id = users.user_id AND fb.brand_id = ?", brandId).
		Where("fp.product_id = ? AND fb.user_id IS NULL", productID).
		Pluck("users.nf_token", &tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func handleBatchResponse(batchResponse *messaging.BatchResponse) {
	log.Println("Batch response:")
	if batchResponse.SuccessCount > 0 {
		log.Printf("sent notification to %d users\n", batchResponse.SuccessCount)
	}
	if batchResponse.FailureCount > 0 {
		log.Printf("failed to send notification to %d users\n", batchResponse.FailureCount)
	}
	for _, result := range batchResponse.Responses {
		if result.Error != nil {
			log.Printf("error sending notification %s: %v\n", result.MessageID, result.Error)
		}
	}
	log.Println("--------------------------------")
}
