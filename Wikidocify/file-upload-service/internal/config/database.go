// config/database.go
// This file contains the configuration for the database
// It defines the InitDB function
// It uses the GORM library to interact with the database

package config

import (
	"fmt"
	"log"
	"os"

	"wikidocify/file-upload-service/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection and performs auto-migration.
func InitDB() {
	log.Println("[DATABASE] Starting database initialization...")

	// Load environment variables from .env file (optional, for local dev)
	_ = godotenv.Load()

	// Build DSN from environment variables with sensible defaults
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "doc_db_admin"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "SecurePass889"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "documents_db"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	log.Printf("[DATABASE] Connecting to database: host=%s, user=%s, dbname=%s, port=%s, sslmode=%s",
		host, user, dbname, port, sslmode)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	log.Println("[DATABASE] Attempting database connection...")
	var err error
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

// GetDB returns the database connection instance.
func GetDB() *gorm.DB {
	if DB == nil {
		log.Println("[DATABASE] Warning: Database connection is nil")
	}
	return DB
}
