// config/database.go
// This file contains the configuration for the database
// It defines the InitDB function
// It uses the GORM library to interact with the database

package config

import (
	"fmt"
	"log"
	"os"

	"github.com/atheeralattar/pbl-week2/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	log.Println("[DATABASE] Starting database initialization...")

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("[DATABASE] Warning: Could not load .env file: %v", err)
	}

	// Build DSN from environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	log.Printf("[DATABASE] Connecting to database: host=%s, user=%s, dbname=%s, port=%s, sslmode=%s",
		host, user, dbname, port, sslmode)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	log.Println("[DATABASE] Attempting database connection...")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[DATABASE] Failed to connect to database: %v", err)
	}
	log.Println("[DATABASE] Database connection established successfully!")

	// Auto migrate the Document model
	log.Println("[DATABASE] Starting database migration...")
	err = DB.AutoMigrate(&models.Document{})
	if err != nil {
		log.Fatalf("[DATABASE] Failed to migrate database: %v", err)
	}
	log.Println("[DATABASE] Database migration completed successfully!")
	log.Println("[DATABASE] Database initialization complete")
}

func GetDB() *gorm.DB {
	if DB == nil {
		log.Println("[DATABASE] Warning: Database connection is nil")
	}
	return DB
}
