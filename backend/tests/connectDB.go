package tests

import (
	// "errors"
	"myapp/database"
	"myapp/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//เชื่อมต่อ database เข้า memory  
func SetupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{},&models.SystemLog{})
	database.DB = db
	return db
}
