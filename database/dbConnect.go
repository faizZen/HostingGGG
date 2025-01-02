package database

import (
	"ggg/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	dsn := os.Getenv("host=pg-37d27cf3-gocrudffff.b.aivencloud.com user=avnadmin password=AVNS_NzwxDhqUBLgMAO3EPhE dbname=defaultdb port=10146 sslmode=require")
	// postgres://avnadmin:AVNS_NzwxDhqUBLgMAO3EPhE@pg-37d27cf3-gocrudffff.b.aivencloud.com:10146/defaultdb?sslmode=require
	// Connect to the database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Connected to PostgreSQL database!")

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
