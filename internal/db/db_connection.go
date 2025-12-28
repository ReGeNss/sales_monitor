package db

import (
	"fmt"
	"log"
	"os"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	user := os.Getenv("DB_USER_NAME")
	password := os.Getenv("DB_USER_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DATABASE_NAME")

	if user == "" || password == "" || host == "" || port == "" || database == "" {
		log.Fatalf("Missing required database environment variables")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, port, database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")

	DB = db
}

func GetDB() *gorm.DB {
	return DB
}