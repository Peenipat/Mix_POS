package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)
// const (
// 	host = "localhost"
// 	port = 5432
// 	user = "myuser"
// 	password = "mypassword"
// 	dbname = "mydatabase"
// )

// var DB *gorm.DB

// func ConnectDB(){
// 	dsn := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
// 	db ,err  := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil{
// 		panic("failed to connect to database")
// 	} 

var DB *gorm.DB

func ConnectDB() {
	// Load .env file
	if loadErr := godotenv.Load(".env"); loadErr != nil {
		log.Println("⚠️ .env not found or not loaded:", loadErr)
	}

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")

		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}

		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=public",
			host, port, user, password, dbname, sslmode,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ failed to connect to database: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL")
	DB = db
}
