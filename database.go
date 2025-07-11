package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"personal_site/models"
)

var DB *gorm.DB

func InitDB() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		return fmt.Errorf("DATABASE_DSN not set in .env")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	DB = db

	if err := DB.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("auto migrate failed: %v", err)
	}

	return nil
}
