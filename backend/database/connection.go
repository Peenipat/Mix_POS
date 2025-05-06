package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// 1. ถ้ามี DATABASE_URL อยู่แล้ว ให้ใช้เลย
	dsn := os.Getenv("DATABASE_URL")

	// 2. ถ้าไม่มีค่อยประกอบจากตัวแปร DB_HOST, DB_PORT, …
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

		// กำหนด DSN ตามรูปแบบที่ GORM/pg driver รองรับ
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=public",
			host, port, user, password, dbname, sslmode,
		)
	}

	// 3. เปิด connection
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ failed to connect to database: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL")
	DB = db
}
